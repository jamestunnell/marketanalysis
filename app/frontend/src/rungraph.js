import van from "vanjs-core"
import Datepicker from 'flowbite-datepicker/Datepicker';
import convertCSVToArray from 'convert-csv-to-array'

import { AppErrorAlert} from './apperror.js';
import { PostJSON } from './backend.js';
import { Button, ButtonCancel } from "./buttons.js";
// import { ChartModal, MakeChartConfig } from "./chart.js";
import { Download } from "./download.js";
import { IconPlay } from './icons.js';
import { ModalBackground, ModalForeground } from './modal.js';

const {div, input, p, label} = van.tags

function makeCSVChartData(text) {
    const records = convertCSVToArray(text, {header: true, separator: ",", type: "object"})

    console.log(`parsed ${records.length} CSV records`)
}

function makeNDJSONChartData(text) {

}

const runGraph = ({id, date, symbol, format}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/${id}/run-day`
        const object = {date, symbol, format}
        const options = {accept: 'application/x-ndjson'}

        console.log("running graph", object)

        PostJSON({route, object, options}).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to run graph", appErr);
    
                    reject(appErr);    
                })
            }

            resp.text().then(text => resolve(text))
        }).catch(err => {
            console.log("failed to make run graph request", err)
            
            reject({
                title: "Action Failed",
                message: "failed to make run graph request",
                details: [err.message],
            })
        });
    });
}

const RunGraph = (graph) => {
    const closed = van.state(false)
    const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring"
    const today = new Date()
    const date = van.state("")
    const symbol = van.state("")
    const format = van.state("ndjson")
    const formatNDJSON = input({
        id: "runFormatNDJSON",
        class: inputClass,
        type: "radio",
        value: "ndjson",
        name: "selectFormat",
        checked: true,
        onchange: (e) => {
            format.val = "ndjson"

            console.log("changed result format to NDJSON")
        },
    });
    const formatCSV = input({
        id: "runFormatCSV",
        class: inputClass,
        type: "radio",
        value: "csv",
        name: "selectFormat",
        onchange: (e) => {
            format.val = "csv"

            console.log("changed result format to CSV")
        },
    });
    const dateInput = input({
        id: "runDate",
        class: inputClass,
        type: "text",
        placeholder: 'Select date',
    });

    dateInput.addEventListener('changeDate', (e) => {
        date.val = dateInput.value
    })

    const modal = ModalBackground(
        ModalForeground(
            {id: "foreground"},
            div(
                {id: "modalContent", class: "flex flex-col rounded-md space-y-4"},
                p({class: "text-lg font-medium font-bold text-center"}, "Run Graph"),
                div(
                    {class: "flex flex-col"},
                    p("Result Format"),
                    label(
                        {for: "runFormatCSV"},
                        div(
                            {class: "flex flex-row"},
                            formatCSV,
                            "CSV",
                        ),
                    ),
                    label(
                        {for: "runFormatNDJSON"},
                        div(
                            {class: "flex flex-row"},
                            formatNDJSON,
                            "NDJSON",
                        ),
                    ),
                ),
                div(
                    {class: "flex flex-col"},
                    label({for: "runDate"}, "Date"),
                    dateInput,
                ),
                div(
                    {class: "flex flex-col"},
                    label({for: "runSymbol"}, "Symbol"),
                    input({
                        id: "runSymbol",
                        class: inputClass,
                        type: "text",
                        placeholder: 'Security symbol (SPY, QQQ, etc.)',
                        oninput: e => symbol.val = e.target.value,
                    }),                        
                ),
                div(
                    {class:"mt-4 flex justify-center"},
                    ButtonCancel({onclick: ()=> closed.val = true, child: "Cancel"}),
                    Button({
                        disabled: van.derive(() => date.val.length == 0),
                        child: [IconPlay, "Run"],
                        onclick: () => {
                            if (date.val === "") {
                                AppErrorAlert({
                                    title: "Invalid Input",
                                    message: "date is empty",
                                    details: [],
                                })

                                return
                            }

                            runGraph({
                                id: graph.id,
                                date: date.val,
                                symbol: symbol.val,
                                format: format.val,
                            }).then(text => {
                                console.log("run graph succeeded")

                                const basename = `${graph.name}_${date.val}`
                                
                                let chartData;
                                
                                switch (format.val) {
                                case "csv":
                                    Download({
                                        filename: basename + ".csv",
                                        blob: new Blob([text], {type: 'text/csv'}),
                                    })

                                    chartData = makeCSVChartData(text)

                                    break
                                case "ndjson":
                                    Download({
                                        filename: basename + ".ndjson",
                                        blob: new Blob([text], {type: 'application/x-ndjson'}),
                                    })

                                    chartData = makeNDJSONChartData(text)

                                    break
                                }

                                // ChartModal(MakeChartConfig(chartData))

                                closed.val = true;
                            }).catch(appErr => {
                                AppErrorAlert(appErr)
                            })
                        },
                    }),
                ),
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

export {RunGraph};