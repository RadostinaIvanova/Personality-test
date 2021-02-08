package main
import(
	"fmt"
	"bufio"
	"net"
//	"os"
	"log"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell"
	"strings"
	"strconv"

	
)
func extractInfo(option string, clientReader bufio.Reader, clientWriter bufio.Writer) {
	clientWriter.WriteString(option + "\n");
	clientWriter.Flush();
	info, err2 := clientReader.ReadString('\n')
	if err2!= nil{
		log.Println(err2.Error())
	}
	fmt.Println(info)
}

func options(app *tview.Application,clientReader bufio.Reader, clientWriter bufio.Writer) string{
	list := tview.NewList()
	var option string
	des := "Description"
	work := "Work"
	soc := "Socialskills"
	rom := "RomanticRelationships"
	list.AddItem(des, "", 'a', func(){ extractInfo(des,clientReader,clientWriter) })
	list.AddItem("Social skills", "", 'b',func(){ extractInfo(soc,clientReader,clientWriter) })
	list.AddItem(work, "", 'c', func(){extractInfo(work,clientReader,clientWriter)})
	list.AddItem("Romantic relationships", "", 'd',func(){ extractInfo(rom,clientReader,clientWriter)})
	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})
	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
	panic(err)
	}
	return option
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

func setButton(app *tview.Application, button *tview.Button){
	button.SetBorder(true).SetRect(0, 0, 50, 3)
	button.SetBackgroundColorActivated(tcell.ColorDefault)
	button.SetLabelColorActivated(tcell.ColorFuchsia)
	if err := app.SetRoot(button, false).SetFocus(button).Run(); err != nil {
		panic(err)
	}
}

func main(){
	conn,err := net.Dial("tcp", "localhost:9000")
	if err != nil{
		log.Fatal(err)
	}
	defer conn.Close()
	clientReader := bufio.NewReaderSize(conn,5000)
	clientWriter := bufio.NewWriterSize(conn,5000)
	
	app := tview.NewApplication()
	button := tview.NewButton("Hit enter to start the TEST").SetSelectedFunc(func() {
		app.Stop()
	})
	setButton(app, button)
	received, err:= clientReader.ReadString('\n')
	if err != nil{
		fmt.Println(err)
	}
	received = strings.TrimSuffix(received, "\n")
	numOfQuestions,_ := strconv.Atoi(received)
	for i := 0; i < numOfQuestions; i++ {
		received, err:= clientReader.ReadString('?')
		if err != nil{
			fmt.Println(err)
		}
		text := getInputFieldText(app,received)
		clientWriter.WriteString(text + "\n")
		clientWriter.Flush()
	}
	classificationResult, err:= clientReader.ReadString('\n')
		if err != nil{
			fmt.Println(err)
		}
	classificationResult = strings.TrimSuffix(classificationResult, "\n")
	fmt.Println(classificationResult)
	var option string
	for ; option!= "Exit";{
		option = options(app, *clientReader,*clientWriter)
	}
}

