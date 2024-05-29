import van from "vanjs-core"

import { Button, ButtonCancel } from "./buttons";
import { ButtonGroup } from "./buttongroup";
import { ModalBackground } from "./modal";
import capitalize from "./capitalize";
import { Table, TableRow } from "./table";

const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

const { div, input, label, p, tbody} = van.tags

const RecordingRow = ({name, flag}) => {
    const props = {
        id: name,
        class: inputClass,
        type: "checkbox",
        onchange: e => flag.val = e.target.checked,
        checked: flag.val,
    }
    const rowItems = [
        label({for: name}, name),
        input(props, capitalize(name)),
    ]

    return TableRow(rowItems)
}

const EditRecordingModal = ({outputNames, recording, onComplete}) => {
    const closed = van.state(false)
    
    outputNames.sort()

    const flags = outputNames.map(name => {
        return van.state(recording.indexOf(name) >= 0)
    })

    const recordingTableBody = tbody({class:"table-auto"}); 
    const reccordingTable = Table({
        columnNames: ["Output", "Record?"],
        tableBody: recordingTableBody,
    })
    const rows = outputNames.map((name, i) => {
        return RecordingRow({name, flag: flags[i]})
    })

    van.add(recordingTableBody, rows)

    const cancelBtn = ButtonCancel({
        child: "Cancel",
        onclick: () => closed.val = true,
    })
    const okBtn = Button({
        child: "OK",
        onclick: () => {
            const nowRecording = outputNames.filter((name, i) => flags[i].val)

            onComplete(nowRecording)

            closed.val = true
        },
    })
    const buttons = ButtonGroup({buttons: [cancelBtn, okBtn], moreClass: "self-end"})
    const modal = ModalBackground(
        div(
            {id: "foreground", class: "flex flex-col block p-16 rounded-lg bg-white min-w-[25%] max-w-[50%]"},
            p({class: "text-lg font-medium font-bold text-center"}, "Edit Recording"),
            reccordingTable,
            buttons,
        )
    )

    van.add(document.body, () => closed.val ? null : modal);
}

export { EditRecordingModal }