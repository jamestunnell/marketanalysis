import van from "vanjs-core"

import { AppErrorAlert} from './apperror.js';
import { Button, ButtonCancel } from "./buttons.js";
import { IconPlay } from './icons.js';
import { DownloadCSV } from "./download.js";
import { PostJSON } from './backend.js';
import { ModalBackground, ModalForeground } from './modal.js';
import Datepicker from 'flowbite-datepicker/Datepicker';

const {div, input, p, label} = van.tags

const runGraph = ({id, date, symbol}) => {
    return new Promise((resolve, reject) => {
        console.log("running graph on date %s", date)

        PostJSON({
            route:`/graphs/${id}/run-day`,
            object: {date: date, symbol: symbol},
            options: {accept: 'text/csv'}
        }).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to run graph", appErr);
    
                    reject(appErr);    
                })
            }

            // expect text/CSV
            resp.text().then(csvData => resolve(csvData))
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
    const closed = van.state(false);
    const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";
    const today = new Date();
    const date = van.state("");
    const symbol = van.state("");

    const dateInput = input({
        id: "runDate",
        class: inputClass,
        type: "text",
        placeholder: 'Select date',
        oninput: (e) => console.log("run date input: %s", e.target.value),
        onchange: (e) => console.log("run date change: %s", e.target.value),
        onselect: (e) => console.log("run date select: %s", e.target.value),
    });

    dateInput.addEventListener('change', function () {
        console.log("changed input value %s", dateInput.val)
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
                label({for: "runDate"}, "Date"),
                dateInput,
                label({for: "runSymbol"}, "Symbol"),
                input({
                    id: "runSymbol",
                    class: inputClass,
                    type: "text",
                    placeholder: 'Security symbol (SPY, QQQ, etc.)',
                    oninput: e => symbol.val = e.target.value,
                }),
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
                            }).then(csvData => {
                                const csvName = `${graph.name}_${date.val}.csv`
        
                                console.log("run graph succeeded, downloading CSV %s", csvName)

                                DownloadCSV({name: csvName, data: csvData})                    

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