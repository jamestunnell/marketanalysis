import van from "vanjs-core"

import { Table, TableRow } from "../table"

const { input, option, select, tbody} = van.tags

const inputClass = "block border border-gray-200 rounded-md focus:border-gray-500 focus:outline-none focus:ring";

const InputSourcesTable = ({inputs, inputSources, possibleSources}) => {
    const names = inputs.map(input => input.name).sort()
    const rows = names.map(name => {
        const selectedSource = inputSources[name]
        const opts = [ option({value: "", selected: selectedSource.val === ""}, "") ]

        possibleSources.forEach(source => {
            const opt = option({value: source, selected: selectedSource.val === source}, source)

            opts.push(opt)
        })

        const selectSource = select(
            {
                class: inputClass,
                oninput: (e) => selectedSource.val = e.target.value,
            },
            opts,
        )

        return TableRow([name, input.type, selectSource])
    })

    return Table({
        columnNames: ["Name", "Type", "Connected Source"],
        tableBody: tbody({class:"table-auto"}, rows),
    })
}

export {InputSourcesTable}