import van from "vanjs-core"
import {Modal} from "vanjs-ui"

import { Button, ButtonCancel, ButtonIcon, ButtonIconTooltip } from './buttons.js'
import { IconCheck, IconDelete, IconError } from './icons.js'
import { TableRow } from './table.js'
import { validateParamVal } from "./paramvals.js"

const {datalist, div, input, label, option, p} = van.tags

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
        const sourceParts = connection.source.split(".")
        const targetParts = connection.target.split(".")

        const sourceBlock = sourceParts[0]
        const sourceOutput = (sourceParts.length > 1) ? sourceParts[1] : ""
        const targetBlock = targetParts[0]
        const targetInput = (targetParts.length > 1) ? targetParts[1] : ""

        console.log(`creating connection row`, {sourceBlock, sourceOutput, targetBlock, targetInput})

        this.id = id
        this.parent = parent
        this.sourceBlock = van.state(sourceBlock)
        this.sourceOutput = van.state(sourceOutput)
        this.targetBlock = van.state(targetBlock)
        this.targetInput = van.state(targetInput)
        this.deleted = van.state(false)

        this.sourceBlocksDatalist = datalist({id:`sourceBlocks-${id}`})
        this.sourceOutputsDatalist = datalist({id:`sourceOutputs-${id}`})
        this.targetBlocksDatalist = datalist({id:`targetBlocks-${id}`})
        this.targetInputsDatalist = datalist({id:`targetInputs-${id}`})

        this.updateSourceDatalistOptions(sourceBlock, sourceOutput)
        this.updateTargetDatalistOptions(targetBlock, targetInput)
    }

    onBlockNameChange() {
        this.updateSourceDatalistOptions(this.sourceBlock.val, this.sourceOutput.val)
        this.updateTargetDatalistOptions(this.targetBlock.val, this.targetInput.val)
    }

    updateSourceDatalistOptions(blockName, outName) {
        this.sourceBlocksDatalist.replaceChildren(
            ...this.parent.blockNames().map(name => {
                return option({value: name, selected: (name === blockName)}, name)
            })
        )

        const sourceInfo = this.parent.findBlockInfo(blockName)
        const sourceOutputOpts = []
        
        if (sourceInfo) {
            sourceInfo.outputs.forEach(out => {
                const opt = option({
                    value: out.name,
                    selected: (out.name === outName),
                }, out.name)

                sourceOutputOpts.push(opt)
            })
        }

        this.sourceOutputsDatalist.replaceChildren(...sourceOutputOpts)
    }

    updateTargetDatalistOptions(blockName, inName) {
        this.targetBlocksDatalist.replaceChildren(
            ...this.parent.blockNames().map(name => {
                return option({value: name, selected: (name === blockName)}, name)
            })
        )

        const targetInfo = this.parent.findBlockInfo(blockName)
        const targetInputOpts = []
        
        if (targetInfo) {
            targetInfo.inputs.forEach(input => {
                const opt = option({
                    value: input.name,
                    selected: (input.name === inName),
                }, input.name)

                targetInputOpts.push(opt)
            })
        }

        this.targetInputsDatalist.replaceChildren(...targetInputOpts)
    }

    makeConnection() {
        return {
            source: `${this.sourceBlock.val}.${this.sourceOutput.val}`,
            target: `${this.targetBlock.val}.${this.targetInput.val}`,
        }
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
            div(
                {class:"text-container"},
                input({
                    class: inputClass,
                    type: "text",
                    list: this.sourceBlocksDatalist.getAttribute('id'),
                    value: this.sourceBlock.val,
                    placeholder: "<source block>",
                    oninput: e => {
                        const sourceBlock = e.target.value
                        
                        this.sourceBlock.val = sourceBlock
                        
                        this.updateSourceDatalistOptions(sourceBlock, this.sourceOutput.val)
                    },
                }),
                this.sourceBlocksDatalist,
            ),
            div(
                {class:"text-container"},
                input({
                    class: inputClass,
                    type: "text",
                    list: this.sourceOutputsDatalist.getAttribute('id'),
                    value: this.sourceOutput.val,
                    placeholder: "<block output>",
                    oninput: e => this.sourceOutput.val = e.target.value,
                }),
                this.sourceOutputsDatalist,
            ),
            div(
                {class:"text-container"},
                input({
                    class: inputClass,
                    type: "text",
                    list: this.targetBlocksDatalist.getAttribute('id'),
                    value: this.targetBlock.val,
                    placeholder: "<target block>",
                    oninput: e => {
                        const targetBlock = e.target.value
                        
                        this.targetBlock.val = targetBlock
                        
                        this.updateTargetDatalistOptions(targetBlock, this.targetInput.val)
                    },
                }),
                this.targetBlocksDatalist,
            ),
            div(
                {class:"text-container"},
                input({
                    class: inputClass,
                    type: "text",
                    list: this.targetInputsDatalist.getAttribute('id'),
                    value: this.targetInput.val,
                    placeholder: "<block input>",
                    oninput: e => this.targetInput.val = e.target.value,
                }),
                this.targetInputsDatalist,
            ),
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