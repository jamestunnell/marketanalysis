import van from 'vanjs-core'
import Datepicker from 'flowbite-datepicker/Datepicker'

import { AppErrorAlert} from './apperror.js'
import { PostJSON } from './backend.js'
import { Button, ButtonIcon } from "./buttons.js"
import { DownloadJSON } from "./download.js"
import { IconDownload, IconClose, IconPlay, IconPlot } from './icons.js'
import { ModalBackground } from './modal.js'
import { PlotRecordingModal } from './plot.js'

const {div, input, p, label, option, select} = van.tags

const evalGraph = ({id, date, symbol, source, predictor, window}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/${id}/eval`
        const object = {type: "slope", date, symbol, source, predictor, window}
        const options = {accept: 'application/json'}

        console.log("evaluating graph", object)

        PostJSON({route, object, options}).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to evaluate graph", appErr);
    
                    reject(appErr);    
                })
            }

            resp.json().then(obj => resolve(obj))
        }).catch(err => {
            console.log("failed to send eval graph request", err)
            
            reject({
                title: "Action Failed",
                message: "failed to send eval graph request",
                details: [err.message],
            })
        });
    });
}

const WINDOW_MIN = 3
const WINDOW_MAX = 100

const EvalGraph = (graph, infoByType) => {
    const closed = van.state(false)
    const completed = van.state(false)
    const recording = van.state({})
    const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring"
    const today = new Date()
    const symbol = van.state("")
    const date = van.state("")
    const window = van.state(WINDOW_MIN)
    const source = van.state("")
    const predictor = van.state("")
    const dateInput = input({
        id: "evalDate",
        class: inputClass,
        type: "text",
        placeholder: 'Select date',
    });

    dateInput.addEventListener('changeDate', (e) => {
        date.val = dateInput.value
    })

    const onRun = () => {
        completed.val = false

        evalGraph({
            id: graph.id,
            symbol: symbol.val,
            date: date.val,
            window: window.val,
            source: source.val,
            predictor: predictor.val,
        }).then(obj => {
            console.log("eval graph succeeded")

            recording.val = obj
            completed.val = true
        }).catch(appErr => {
            AppErrorAlert(appErr)
        })
    }

    const closeBtn = ButtonIcon({icon: IconClose(), onclick: ()=> closed.val = true})
    const runBtn = Button({
        disabled: van.derive(() => {
            return (
                (date.val.length === 0) || 
                (symbol.val.length === 0) ||
                (window.val < WINDOW_MIN) || 
                (window.val > WINDOW_MAX) ||
                (source.val.length === 0) || 
                (predictor.val.length === 0)
            )
        }),
        child: [IconPlay, " Run"],
        onclick: onRun,
    })
    const plotBtn = ButtonIcon({
        icon: IconPlot(),
        onclick: ()=> PlotRecordingModal(recording.val),
    })
    const downloadBtn = ButtonIcon({
        icon: IconDownload(),
        onclick: ()=> {
            DownloadJSON({
                filename: `${graph.name}_${symbol.val}_${date.val}_eval.json`,
                object: recording.val,
            })
        },
    })
    const resultsArea = div(
        {
            class: van.derive(() => `flex flex-row justify-center ${!completed.val ? "hidden" : ""}`),
        },
        plotBtn,
        downloadBtn,
    )

    closeBtn.classList.add("self-end")
    runBtn.classList.add("self-center")

    const blockOuts = []
    graph.blocks.forEach(blk => {
        infoByType[blk.type].outputs.forEach(out => {
            blockOuts.push(blk.name + "." + out.name)
        })
    })
    
    console.log(`made ${blockOuts.length} block outputs`, blockOuts)

    const sourceBlockOutOpts = [ option({value:"", selected: true}, "") ].concat(
        blockOuts.map(blkOut => option({value: blkOut}, blkOut))
    )
    const predBlockOutOpts = [ option({value:"", selected: true}, "") ].concat(
        blockOuts.map(blkOut => option({value: blkOut}, blkOut))
    )

    const modal = ModalBackground(
        div(
            {id: "foreground", class: "block p-16 rounded-lg bg-white min-w-[25%] max-w-[25%]"},
            div(
                {id: "modalContent", class: "flex flex-col rounded-md space-y-4"},
                closeBtn,
                p({class: "text-lg font-medium font-bold text-center"}, "Evaluate Graph"),
                div(
                    {class: "flex flex-col"},
                    label({for: "evalDate"}, "Date"),
                    dateInput,
                    label({for: "window"}, "Window Size"),
                    input({
                        id: "window",
                        type: "number",
                        class: inputClass,
                        value: window.val,
                        min: WINDOW_MIN,
                        max: WINDOW_MAX,
                        step: 1,
                        onchange: e => window.val = Number(e.target.value),
                    }),
                    label({for: "source"}, "Source"),
                    select({
                        id: "source",
                        class: inputClass,
                        oninput: e => source.val = e.target.value,
                    }, sourceBlockOutOpts),
                    label({for: "predictor"}, "Predictor"),
                    select({
                        id: "predictor",
                        class: inputClass,
                        oninput: e => predictor.val = e.target.value,
                    }, predBlockOutOpts),
                ),
                div(
                    {class: "flex flex-col"},
                    label({for: "runSymbol"}, "Symbol"),
                    input({
                        id: "runSymbol",
                        class: inputClass,
                        type: "text",
                        placeholder: 'Symbol (SPY, QQQ, etc.)',
                        onchange: e => symbol.val = e.target.value,
                    }),                        
                ),
                runBtn,
                resultsArea,
            ),
        ),
    )
  
    van.add(document.body, () => closed.val ? null : modal);

    const datePickerOpts = {
        autohide: true,
        container: "#foreground",
        daysOfWeekDisabled: [0, 6], // disable saturday and sunday
        format: "yyyy-mm-dd",
        maxDate: today,
        todayHighlight: true,
    }
    const datepicker = new Datepicker(dateInput, datePickerOpts)
}

export { EvalGraph };