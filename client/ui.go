package main
import(
	"fmt"
	"bufio"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
	"strings"
	"log"
)

const startTest string = "Hit enter to start the TEST"

func textView(text string){
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})
		go func() {
			for _, word := range strings.Split(text, " ") {
				fmt.Fprintf(textView, "%s ", word)
			}
		}()
	textView.SetDoneFunc(func(key tcell.Key) {
			app.Stop()
	})
	textView.SetBorder(true)
	if err := app.SetRoot(textView, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}


func setButton(app *tview.Application, button *tview.Button){
	button.SetBorder(true).SetRect(0, 0, 50, 3)
	button.SetBackgroundColorActivated(tcell.ColorDefault)
	button.SetLabelColorActivated(tcell.ColorFuchsia)
	if err := app.SetRoot(button, false).SetFocus(button).Run(); err != nil {
		panic(err)
	}
}

func setStartTestButton(text string, app *tview.Application){
	button := tview.NewButton(text).SetSelectedFunc(func() {
		app.Stop()
	})
	setButton(app, button)
}
func options(clientReader bufio.Reader, clientWriter bufio.Writer){
	app := tview.NewApplication()
	list := tview.NewList()
	des := "Description"
	work := "Work"
	soc := "Socialskills"
	rom := "RomanticRelationships"
	list.AddItem(des, "", 'a', func(){  app.Stop()  
		receiveInfo(des,clientReader,clientWriter) 
	app.Stop()})
	list.AddItem("Social skills", "", 'b',func(){ app.Stop()
		 receiveInfo(soc,clientReader,clientWriter)
		})
	list.AddItem(work, "", 'c', func(){app.Stop()
		 receiveInfo(work,clientReader,clientWriter)
		})
	list.AddItem("Romantic relationships", "", 'd',func(){app.Stop() 
		 receiveInfo(rom,clientReader,clientWriter)
		})
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
		clientWriter.WriteString("Quit" + "\n");
		clientWriter.Flush();
		
	})
	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}

func getInputFieldText(app *tview.Application, received string) string{
	inputField := tview.NewInputField()
	inputField.SetLabel(received)
	inputField.SetFormAttributes(len(inputField.GetLabel()),tcell.ColorLime, tcell.ColorDefault , tcell.ColorFuchsia, tcell.ColorDefault)
	inputField.SetAcceptanceFunc(tview.InputFieldMaxLength(64))
		inputField.SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})
		if err := app.SetRoot(inputField, true).SetFocus(inputField).Run(); err != nil {
			panic(err)
		}
	return inputField.GetText()
}

func modalClassification(app *tview.Application , text string){
	modal := tview.NewModal()
	modal.SetBackgroundColor(tcell.ColorLime)
	modal.SetTextColor(tcell.ColorFuchsia)
	modal.AddButtons([]string {""})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				app.Stop()
	})
	modal.SetText("Your personality type is " + text)
	if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
		panic(err)
	}
}

func receiveInfo(option string, clientReader bufio.Reader, clientWriter bufio.Writer) {
	clientWriter.WriteString(option + "\n");
	clientWriter.Flush();
	info, err2 := clientReader.ReadString('\n')
	if err2!= nil{
		log.Println(err2.Error())
	}
	textView(info)
	options(clientReader,clientWriter)
}