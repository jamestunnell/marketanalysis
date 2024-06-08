import van from "vanjs-core"

import { Table, TableRow } from "./table"

const { div, input, label, option, p, select, tbody} = van.tags

const inputClass = "block border border-gray-200 rounded-md focus:border-gray-500 focus:outline-none focus:ring";

const RecordedOutputsTable = ({outputs, recordedOutputs}) => {
    const names = outputs.map(output => output.name).sort()
    const rows = names.map(name => {
        const recorded = recordedOutputs[name]
        const checkbox = input({
            class: inputClass,
            type: "checkbox",
            onchange: e => recorded.val = e.target.checked,
            checked: recorded.val,
        })

        return TableRow([name, checkbox])
    })

    return Table({
        columnNames: ["Name", "Recorded"],
        tableBody: tbody({class:"table-auto"}, rows),
    })
}

export { RecordedOutputsTable }