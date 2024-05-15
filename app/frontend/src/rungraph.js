import van from "vanjs-core"

import { DoAppErrorModal} from './apperror.js';
import { Button, ButtonCancel } from "./buttons.js";
import { IconPlay } from './icons.js';
import { DownloadCSV } from "./download.js";
import { Post } from './backend.js';
import { ModalBackground, ModalForeground } from './modal.js';
import Datepicker from '@themesberg/tailwind-datepicker/Datepicker';

const {div, input, p, label} = van.tags

const RunGraph = (graph) => {
    const closed = van.state(false);
    const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";
    const today = new Date();
    const date = van.state(today.toString());


    const dateInput = input({
        id: "dateInput",
        class: inputClass,
        type: "text",
        value: `${today.getFullYear()}-${today.getMonth()}-${today.getDay()}`,
        oninput: e => date.val = e.target.value,
    });
    const modal = ModalBackground(
        ModalForeground(
            {id: "foreground"},
            div(
                {id: "modalContent", class: "flex flex-col rounded-md space-y-4"},
                p({class: "text-lg font-medium font-bold text-center"}, "Run Graph"),
                label({for: "dateInput"}, "Date"),
                dateInput,
                div(
                    {class:"mt-4 flex justify-center"},
                    ButtonCancel({onclick: ()=> closed.val = true, child: "Cancel"}),
                    Button({onclick: ()=> closed.val = true, child: [IconPlay, "Run"]}),
                ),
            ),
        ),
    )

    van.add(document.body, () => closed.val ? null : modal);

    const datePickerOpts = {
        autohide: true,
        daysOfWeekDisabled: [0, 6], // disable saturday and sunday
        daysOfWeekHighlighted: [1,2,3,4,5], // highlight weekdays
        format: "yyyy-mm-dd",
        maxDate: today,
        todayHighlight: true,
    }
    const datepicker = new Datepicker(dateInput, datePickerOpts)

        // onErr: (appErr) => DoAppErrorModal(appErr),
        // onSuccess: (csvData) => {
        //     name = graph.name + "_" + ;

        //     DownloadCSV({name: name, data: csvData})
        // },
}

export {RunGraph};