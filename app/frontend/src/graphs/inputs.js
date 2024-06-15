import van from "vanjs-core"

import { Table, TableRow } from "../table"
import { INPUT_CLASS } from "../input"

const { input, option, select, tbody} = van.tags

const InputsTable = ({inputs, sources, possibleSources}) => {
    const names = inputs.map(input => input.name).sort()
    const rows = names.map(name => {
        const selectedSource = sources[name]
        const opts = [ option({value: "", selected: selectedSource.val === ""}, "") ]

        possibleSources.forEach(source => {
            const opt = option({value: source, selected: selectedSource.val === source}, source)

            opts.push(opt)
        })

        const selectSource = select(
            {
                class: INPUT_CLASS,
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

export {InputsTable}