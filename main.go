package main

import (
	"encoding/json"
	"net/http"

	"github.com/syumai/workers"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type PlotRequest struct {
	Data   []float64 `json:"data"`
	Title  *string   `json:"title,omitempty"`
	XLabel *string   `json:"x_label,omitempty"`
	YLabel *string   `json:"y_label,omitempty"`
}

func main() {
	// helloエンドポイント: テスト用
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Hello!"))
	})

	// /plotエンドポイント: JSON配列を受け取り、グラフをPNGで返す
	http.HandleFunc("/plot", func(w http.ResponseWriter, req *http.Request) {
		// JSONデコード
		var plotReq PlotRequest
		if err := json.NewDecoder(req.Body).Decode(&plotReq); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if len(plotReq.Data) == 0 {
			http.Error(w, "data array is empty", http.StatusBadRequest)
			return
		}

		// グラフ生成
		p := plot.New()

		// タイトルとラベルの設定（指定がある場合のみ）
		if plotReq.Title != nil {
			p.Title.Text = *plotReq.Title
		} else {
			p.Title.Text = "Example Plot"
		}

		if plotReq.XLabel != nil {
			p.X.Label.Text = *plotReq.XLabel
		} else {
			p.X.Label.Text = "X"
		}

		if plotReq.YLabel != nil {
			p.Y.Label.Text = *plotReq.YLabel
		} else {
			p.Y.Label.Text = "Y"
		}

		// dataをプロット用の点列に変換
		pts := make(plotter.XYs, len(plotReq.Data))
		for i, v := range plotReq.Data {
			pts[i].X = float64(i)
			pts[i].Y = v
		}

		line, err := plotter.NewLine(pts)
		if err != nil {
			http.Error(w, "failed to create line plotter", http.StatusInternalServerError)
			return
		}
		p.Add(line)

		// レスポンスをPNGで返す
		w.Header().Set("Content-Type", "image/png")

		// p.Saveではファイルに書き出すため、代わりに直接ResponseWriterに出力するには
		// p.WriterToを使う
		writerTo, err := p.WriterTo(4*vg.Inch, 4*vg.Inch, "png")
		if err != nil {
			http.Error(w, "failed to create writer", http.StatusInternalServerError)
			return
		}
		_, err = writerTo.WriteTo(w)
		if err != nil {
			http.Error(w, "failed to write image", http.StatusInternalServerError)
			return
		}
	})

	workers.Serve(nil) // use http.DefaultServeMux
}
