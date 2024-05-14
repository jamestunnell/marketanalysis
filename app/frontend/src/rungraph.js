import van from "vanjs-core"

import { DoAppErrorModal} from './apperror.js';
import { Button, ButtonCancel } from "./buttons.js";
import { IconPlay } from './icons.js';
import { DownloadCSV } from "./download.js";
import { Post } from './backend.js';
import { Modal } from "vanjs-ui";
import Datepicker from 'vanillajs-datepicker/Datepicker';

const {div, input, p, label} = van.tags


const RunGraph = (graph) => {
    const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";
    const today = new Date();
    const date = van.state(today.toString());

    const dateInput = input({
        id: "date",
        class: inputClass,
        type: "text",
        value: `${today.getFullYear()}-${today.getMonth()}-${today.getDay()}`,
        oninput: e => date.val = e.target.value,
    });

    const datePickerOpts = {
        autohide: true,
        daysOfWeekDisabled: [0, 6], // disable saturday and sunday
        format: "yyyy-mm-dd",
        maxDate: today,
    }
    const datepicker = new Datepicker(dateInput, datePickerOpts)

    const form = div(
        {class: "flex flex-col rounded-md space-y-4"},
        p({class: "text-lg font-medium font-bold text-center"}, "Run Graph"),
        label({for: "date"}, "Date"),
        dateInput,
        div(
            {class:"mt-4 flex justify-center"},
            ButtonCancel({onclick: ()=> closed.val = true, child: "Cancel"}),
            Button({onclick: ()=> closed.val = true, child: [IconPlay, "Run"]}),
        ),
    )

    van.add(document.body, Modal({closed}, form))
        // onErr: (appErr) => DoAppErrorModal(appErr),
        // onSuccess: (csvData) => {
        //     name = graph.name + "_" + ;

        //     DownloadCSV({name: name, data: csvData})
        // },
}

export {RunGraph};