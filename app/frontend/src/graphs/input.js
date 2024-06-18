import van from "vanjs-core"

import Select from '../elements/select.js'
import { Table, TableRow } from '../elements/table.js'

const { option, tbody } = van.tags

class Input {
    constructor({name, source, type}) {
        this.name = name
        this.source = source
        this.type = type
    }

    isSourceEmpty() {
        return this.source.length === 0
    }

    makeConfig() {
        return {name: this.name, source: this.source.val}
    }

    renderTableRow(possibleSources) {
        const opts = [ option({value: "", selected: this.source.val === ""}, "") ]

        possibleSources.forEach(source => {
            const opt = option({value: source, selected: this.source.val === source}, source)

            opts.push(opt)
        })

        const selectSource = Select({
            onchange: (e) => this.source.val = e.target.value,
            options: opts,
        })

        return TableRow([this.name, this.type, selectSource])
    }
}

function InputsTable({inputs, hidden, possibleSources}) {
    return Table({
        hidden,
        columnNames: ["Name", "Type", "Connected Source"],
        tableBody: tbody({class:"table-auto"}, inputs.map(input => input.renderTableRow(possibleSources))),
    })
}

function MakeInputs({infos, configs}) {
    return infos.map(info => {
        const cfg = configs.find(cfg => cfg.name === info.name)
        
        return new Input({
            name: info.name,
            type: info.type,
            source: van.state(cfg ? cfg.source : ""),
        })
    })
}

export {InputsTable, MakeInputs}