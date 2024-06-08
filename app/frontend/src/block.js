import van from "vanjs-core"
import hash from 'object-hash';

import { Button, ButtonCancel, ButtonIcon, ButtonIconTooltip } from './buttons.js';
import { ButtonGroup } from "./buttongroup.js";
import { IconCheck, IconClose, IconCollapsed, IconDelete, IconError, IconExpanded } from "./icons.js";
import { ParamValsTable, validateParamVal } from './paramvals.js'
import { EditRecordingModal } from "./recording.js";
import { ModalBackground, ModalForeground } from "./modal.js";
import { TableRow } from './table.js';
import { truncateString } from "./truncatestring.js";

const {div, input, li, ul, option, p, select, span} = van.tags

const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-gray-500 focus:outline-none focus:ring";

function validateBlock({block, info, otherNames}) {
    console.log("validating block", block)

    if (block.name.length === 0) {
        return new Error("Name is empty")
    }

    if (otherNames.indexOf(block.name) >= 0) {
        return new Error(`Name '${block.name}' is not unique`)
    }

    const paramErrs = Object.entries(block.paramVals).map(([name, val]) => {
        const param = info.params.find(p => p.name === name)
        if (!param) {
            // params do not have to be set (default will be used)
            return null
        }

        const err = validateParamVal({param, value: val})
        if (err) {
            return new Error(`param ${param.name} has invalid value ${val}`)
        }
    }).filter(err => err)

    if (paramErrs.length > 0) {
        return paramErrs[0]
    }

    const recordingErrs = block.recording.map(name => {
        const out = info.outputs.find(o => o.name === name)

        return out ? null : new Error(`failed to find output ${name} marked for recording`)
    }).filter(err => err)

    if (recordingErrs.length > 0) {
        return recordingErrs[0]
    }

    return null
}

class BlockItem {
    constructor({id, block, info, parent}) {
        this.id = id
        this.info = info
        this.parent = parent
        this.type = block.type

        this.block = van.state(block)
        this.name = van.derive(() => this.block.val.name)
    }

    getName() {
        return this.name.val
    }
    
    makeBlock() {
        return this.block.val
    }

    delete() {
        this.parent.deleteBlock(this.id)
    }

    renderButton() {
        return Button({child: this.name, onclick: () => this.editModal()})
    }

    editModal() {
        const blockBefore = this.block.val
        const connectionsBefore = this.parent.connections
        const closed = van.state(false)
        const onComplete = ({block, connections}) => {
            let changed = false

            if (hash(block) !== hash(blockBefore)) {
                this.block.val = block
                
                changed = true
            }

            if (hash(connections) !== hash(connectionsBefore)) {
                this.parent.connections = connections

                changed = true
            }

            if (changed) {
                this.parent.updateDigest()
            }
    
            closed.val = true
        }
        const onCancel = () => closed.val = true;
        const onDelete  = () => {
            closed.val = true

            this.parent.deleteBlock(this.id)
        }
        const form = EditBlockForm({
            block: blockBefore,
            info: this.info,
            otherNames: this.parent.blockNames().filter(name => name !== blockBefore.name),
            connections: connectionsBefore,
            onComplete, onCancel, onDelete,
        })
        const modal = ModalBackground(
            div(
                {class: "block p-8 rounded-lg bg-white z-11"},
                form,
            ),
        )
    
        van.add(document.body, () => closed.val ? null : modal);
    }
}

const EditBlockForm = ({block, info, otherNames, connections, onComplete, onCancel, onDelete}) => {
    const type = block.type
    const nameWorking = van.state(block.name)
    const paramValsWorking = Object.fromEntries(info.params.map(p => {
        const nonDefaultVal = block.paramVals[p.name]

        return [p.name, van.state(nonDefaultVal ?? p.default)]
    }))
    const paramsCollapsed = van.state(false)
    const inputsCollapsed = van.state(false)
    const outputsCollapsed = van.state(false)
    const modifiedBlock = van.derive(() => {
        const nonDefaultVals = {}

        Object.entries(paramValsWorking).forEach(([name, value]) => {
            if (value.val !== info.params.find(p => p.name === name).default) {
                nonDefaultVals[name] = value.val
            }
        })

        return {name: nameWorking.val, type: block.type, paramVals: nonDefaultVals, recording: block.recording}
    })
    const nameInput = input({
        class: inputClass,
        type: "text",
        value: nameWorking,
        placeholder: "Non-empty, unique",
        oninput: e => nameWorking.val = e.target.value,
    })
    const closeBtn = ButtonIcon({icon: IconClose(), onclick: onCancel})
    const deleteBtn = ButtonIcon({icon: IconDelete(), onclick: onDelete});
    const validateErr = van.derive(() => {
        console.log("using other names", otherNames)
        
        return validateBlock({
            block: modifiedBlock.val,
            info,
            otherNames,
        })
    })
    const statusBtn = ButtonIconTooltip({
        icon: () => validateErr.val ? IconError() : IconCheck(),
        tooltipText: van.derive(() => validateErr.val ? `Block is invalid: ${validateErr.val.message}` : "Block is valid"),
    });
    const ok = Button({
        child: "OK",
        disabled: validateErr,
        onclick: () => onComplete({block: modifiedBlock.val}),
    })
    const cancel = ButtonCancel({child: "Cancel", onclick: onCancel})

    return div(
        {class: "flex flex-col divide-y"},
        div(
            {class: "grid grid-cols-3"},
            deleteBtn, statusBtn, closeBtn
        ),
        div(
            {class: "grid grid-cols-2 space-y-2"},

            span({class: "text-md font-medium font-bold"}, "Name"),
            nameInput,

            span({class: "text-md font-medium font-bold"},"Type"),
            span(type),

            span({class: "mb-2"}),
        ),
        div(
            {class:"flex flex-col"},
            div(
                {class: "flex flex-row items-center"},
                span({class: "text-xl mt-2 font-medium font-bold"}, "Parameters"),
                ButtonIcon({
                    icon: () => paramsCollapsed.val ? IconCollapsed() : IconExpanded(),
                    onclick: () => paramsCollapsed.val = !paramsCollapsed.val,
                })
            ),
            div(
                {hidden: paramsCollapsed},
                ParamValsTable({
                    params: info.params,
                    paramVals: paramValsWorking,
                }),    
            ),
            span({class: "mb-2"}),
        ),
        div(
            {class:"flex flex-col"},
            div(
                {class: "flex flex-row items-center"},
                span({class: "text-xl mt-2 font-medium font-bold"}, "Inputs"),
                ButtonIcon({
                    icon: () => inputsCollapsed.val ? IconCollapsed() : IconExpanded(),
                    onclick: () => inputsCollapsed.val = !inputsCollapsed.val,
                })
            ),
            div(
                {hidden: inputsCollapsed},
                info.inputs.map(input => span(input.name)),
            ),
            span({class: "mb-2"}),
        ),
        div(
            {class:"flex flex-col"},
            div(
                {class: "flex flex-row items-center"},
                span({class: "text-xl mt-2 font-medium font-bold"}, "Outputs"),
                ButtonIcon({
                    icon: () => outputsCollapsed.val ? IconCollapsed() : IconExpanded(),
                    onclick: () => outputsCollapsed.val = !outputsCollapsed.val,
                })
            ),
            div(
                {hidden: outputsCollapsed},
                info.outputs.map(output => span(output.name)),
            ),
            span({class: "mb-2"}),
        ),
        div({class: "flex flex-row-reverse"}, ok, cancel),
    )
}

const AddBlockForm = ({infoByType, blockNames, onComplete, onCancel}) => {
    const types = Object.keys(infoByType)
    const selectedType = van.state(types[0])
    const description = van.derive(() => infoByType[selectedType.val].description)
    const options = types.map((t) => {
        return option({value: t, selected: (t === selectedType.val)}, t);
    })
    const selectType = select(
        {
            id: "type",
            class: inputClass,
            onchange: (e) => selectedType.val = e.target.value,
        },
        options,
    )
    const ok = Button({
        child: "OK",
        onclick: () => {
            const info = infoByType[selectedType.val]
            let name = truncateString(selectedType.val, 3).toLowerCase()

            if (blockNames.indexOf(name) >= 0) {
                let i = 2
                const candidate = () => {return `${name}${i}`}

                while(blockNames.indexOf(candidate()) >= 0) {
                    i++
                }

                name = candidate()
            }

            const block = {type: selectedType.val, name, paramVals: {}, recording: []}
    
            onComplete({info, block})
        },
    })
    const cancel = ButtonCancel({child: "Cancel", onclick: onCancel})

    return div(
        {class: "flex flex-col max-w-250"},
        div(
            {class: "grid grid-cols-2"},
            span({class: "min-w-24 max-w-24"}, p({class: "text-md font-medium font-bold"}, "Type")),
            span({class: "min-w-48 max-w-48"}, selectType),
            span({class: "min-w-24 max-w-24"}, p({class: "text-md font-medium font-bold"}, "Description")),
            span({class: "min-w-48 max-w-48"}, description)
        ),
        div({class:"mt-4 flex flew-row-reverse"}, ok, cancel),
    )
}

const AddBlockModal = ({infoByType, blockNames, handleResult}) => {
    const closed = van.state(false)
    const onComplete = ({block, info}) => {
        handleResult({block, info})

        closed.val = true
    }
    const onCancel = () => closed.val = true;
    const form = AddBlockForm({infoByType, blockNames, onComplete, onCancel})
    const modal = ModalBackground(ModalForeground({}, form))

    van.add(document.body, () => closed.val ? null : modal);
}

export {BlockItem, AddBlockModal};