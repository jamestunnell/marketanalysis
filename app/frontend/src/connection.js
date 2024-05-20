import van from "vanjs-core"
import {Tooltip, Modal} from "vanjs-ui"

import { Button, ButtonCancel, ButtonIcon, ButtonIconTooltip } from './buttons.js'
import { IconCheck, IconDelete, IconError } from './icons.js'
import { TableRow } from './table.js'

const {div, input, label, p} = van.tags

const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

function validateConnection({connection, findBlockInfo}) {
    const srcParts = connection.source.split(".")
    const tgtParts = connection.target.split(".")

    if (srcParts.length !== 2) {
        return new Error(`source ${connection.source} not formatted as <A>.<B>`)
    }

    if (tgtParts.length !== 2) {
        return new Error(`target ${connection.target} not formatted as <A>.<B>`)
    }

    const srcInfo = findBlockInfo(srcParts[0])
    const tgtInfo = findBlockInfo(tgtParts[0])

    if (!srcInfo) {
        return new Error(`source block ${srcParts[0]} not found`)
    }

    if (!tgtInfo) {
        return new Error(`target block ${tgtParts[0]} not found`)
    }

    console.log(`found info for target block ${tgtParts[0]}`)

    if (!srcInfo.outputs.find(o => o.name === srcParts[1])) {
        return new Error(`source block ${srcParts[0]} does not have output ${srcParts[1]}`)
    }
    
    if (!tgtInfo.inputs.find(i => i.name === tgtParts[1])) {
        return new Error(`target block ${tgtParts[0]} does not have input ${tgtParts[1]}`)
    }

    return null
}

class ConnectionRow {
    constructor({id, connection, parent}) {
        this.id = id
        this.parent = parent
        this.source = van.state(connection.source)
        this.target = van.state(connection.target)
        this.deleted = van.state(false)
    }

    makeConnection() {
        return {source: this.source.val, target: this.target.val}
    }

    delete() {
        this.delete.val = true
    }

    render() {
        const deleteBtn = ButtonIcon({
            icon: IconDelete(),
            onclick: () => {
                this.deleted.val = true
    
                this.parent.deleteConnectionRow(this.id)
            },
        });
        const validateErr = van.derive(() => {
            return validateConnection({
                connection: this.makeConnection(),
                findBlockInfo: (name) => this.parent.findBlockInfo(name),
            })
        })
        const statusBtn = ButtonIconTooltip({
            icon: () => validateErr.val ? IconError() : IconCheck(),
            tooltipText: van.derive(() => validateErr.val ? `Connection is invalid: ${validateErr.val.message}` : "Connection is valid"),
        });
    
        const rowItems = [
            input({
                class: inputClass,
                type: "text",
                value: this.source.val,
                placeholder: "<block.output>",
                oninput: e => this.source.val = e.target.value,
            }),
            input({
                class: inputClass,
                type: "text",
                value: this.target.val,
                placeholder: "<block.output>",
                oninput: e => this.target.val = e.target.value,
            }),
            deleteBtn,
            statusBtn
        ];
    
        return () => this.deleted.val ? null : TableRow(rowItems);
    }
}

const ConnectionForm = ({source, target, onOK, onCancel}) => {
    const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

    return div(
        {class: "flex flex-col rounded-md space-y-4"},
        p({class: "text-lg font-medium font-bold text-center"}, "Graph Connection"),
        label({for: "source"}, "Source Address"),
        input({
            id: "source",
            class: inputClass,
            type: "text",
            value: source,
            placeholder: "<block.output>",
            oninput: e => {
                source.val = e.target.value

                this.parent.markChanged()
            },
        }),
        label({for: "target"}, "Target Address"),
        input({
            id: "target",
            class: inputClass,
            type: "text",
            value: target,
            placeholder: "<block.input>",
            oninput: e => {
                target.val = e.target.value

                this.parent.markChanged()
            },
        }),
        div(
            {class:"mt-4 flex justify-center"},
            ButtonCancel({onclick: onCancel, child: "Cancel"}),
            Button({onclick: onOK, child: "OK"}),
        ),
    )
}

const DoConnectionModal = ({connection, handleResult}) => {
    const closed = van.state(false);
    const source = van.state(connection.source);
    const target = van.state(connection.target);

    van.add(
        document.body,
        Modal({closed},
            ConnectionForm({
                source: source,
                target: target,
                onOK: ()=> {
                    handleResult({source: source.val, target: target.val});

                    closed.val = true;
                },
                onCancel: () => {
                    closed.val = true;
                }
            }),
        ),
    );
}

export {ConnectionRow, DoConnectionModal, validateConnection};