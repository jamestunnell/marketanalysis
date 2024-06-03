import van from "vanjs-core"
import hash from 'object-hash';

import { Button, ButtonCancel, ButtonIcon, ButtonIconTooltip } from './buttons.js';
import { ButtonGroup } from "./buttongroup.js";
import { IconCheck, IconEdit, IconDelete, IconError, IconView } from "./icons.js";
import { EditParamValsModal, validateParamVal } from './paramvals.js'
import { EditRecordingModal } from "./recording.js";
import { ModalBackground, ModalForeground } from "./modal.js";
import { TableRow } from './table.js';

const {div, input, label, li, ul, option, p, select, span} = van.tags

const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

function validateBlock({block, info, otherNames}) {
    console.log("validating block", block)

    if (otherNames.indexOf(block.name) >= 0) {
        return new Error(`Name '${block.name}' is not unique`)
    }

    const paramErrs = Object.entries(block.paramVals).map(([name, val]) => {
        const param = info.params.find(p => p.name === name)
        if (!param) {
            // params do not have to be set (default will be used)
            return null
        }

        const err = validateParamVal(param, val)
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

class BlockRow {
    constructor({id, block, info, parent}) {
        this.id = id
        this.info = info
        this.parent = parent
        this.deleted = van.state(false)

        this.type = block.type
        this.name = van.state(block.name)
        this.paramVals = van.state(block.paramVals)
        this.recording = van.state(block.recording)
    }

    getName() {
        return this.name.val
    }
    
    makeBlock() {
        return {
            name: this.name.val,
            type: this.type,
            paramVals: this.paramVals.val,
            recording: this.recording.val,
        }
    }

    delete() {
        this.deleted.val = true
    }

    render() {
        const nameInput = input({
            class: inputClass,
            type: "text",
            value: this.name.val,
            placeholder: "Non-empty, unique",
            oninput: e => {
                this.name.val = e.target.value

                this.parent.onBlockNameChange()
                this.parent.updateDigest()
            },
        })
        const deleteBtn = ButtonIcon({
            icon: IconDelete(),
            // text: "Delete",
            onclick: () => {
                this.deleted.val = true
    
                this.parent.deleteBlockRow(this.id)
            },
        });
        const validateErr = van.derive(() => {
            const otherRows = this.parent.blockRowsWithoutID(this.id)
            const otherNames = otherRows.map(r => r.name.val)
            
            console.log("using other names", otherNames)
            
            return validateBlock({
                block : this.makeBlock(),
                info: this.info,
                otherNames: otherNames,
            })
        })
        const statusBtn = ButtonIconTooltip({
            icon: () => validateErr.val ? IconError() : IconCheck(),
            tooltipText: van.derive(() => validateErr.val ? `Block is invalid: ${validateErr.val.message}` : "Block is valid"),
        });
        const viewParamsBtn = ButtonIconTooltip({
            icon: IconView(),
            // text: "View",
            tooltipText: () => {
                const items = Object.entries(this.paramVals.val).map(([name, val]) => {
                    return li(`${name}: ${val}`)
                })

                return items.length === 0 ? p("All values set to defaults") : ul(items)
            },
        })
        const viewRecordingBtn = ButtonIconTooltip({
            icon: IconView(),
            // text: "View",
            tooltipText: () => {
                const items = this.recording.val.map(name => li(name))

                return items.length === 0 ? p("No outputs set to record") : ul(items)
            },
        })
        const editParamsBtn = ButtonIcon({
            icon: IconEdit(),
            // text: "Edit",
            onclick: () => {
                EditParamValsModal({
                    params: this.info.params,
                    paramVals: this.paramVals.val,
                    onComplete: (paramVals) => {
                        if (hash(paramVals) === hash(this.paramVals.val)) {
                            return
                        }

                        this.paramVals.val = paramVals
                        this.parent.updateDigest()
                    },
                })
            },
        })
        const editRecordingBtn = ButtonIcon({
            icon: IconEdit(),
            // text: "Edit",
            onclick: () => {
                EditRecordingModal({
                    outputNames: this.info.outputs.map(o => o.name),
                    recording: this.recording.val,
                    onComplete: (recording) => {
                        if (hash(recording) === hash(this.recording.val)) {
                            return
                        }

                        this.recording.val = recording
                        this.parent.updateDigest()
                    },
                })
            },
        })

        const paramButtons = ButtonGroup({buttons: [viewParamsBtn, editParamsBtn]})
        const recordingButtons = ButtonGroup({buttons: [viewRecordingBtn, editRecordingBtn]})
        const rowItems = [ nameInput, this.type, paramButtons, recordingButtons, deleteBtn, statusBtn]
    
        return () => this.deleted.val ? null : TableRow(rowItems);
    }
}

const SelectBlockTypeForm = ({types, onComplete, onCancel}) => {
    const selectedType = van.state(types[0])
    const options = types.map((t) => {
        return option({value: t, selected: (t === selectedType.val)}, t);
    })
    const selectType = select(
        { id: "type", class: inputClass, onchange: (e) => selectedType.val = e.target.value },
        options,
    )
    const ok = Button({
        child: "OK",
        onclick: () => onComplete(selectedType.val),
    })
    const cancel = ButtonCancel({child: "Cancel", onclick: onCancel})

    return div(
        {class: "flex flex-col"},
        p({class: "text-lg font-medium font-bold text-center"}, "Block Type"),
        selectType,
        div({class:"mt-4 flex flew-row-reverse"}, ok, cancel),
    )
}

const SelectBlockTypeModal = ({types, handleResult}) => {
    const closed = van.state(false)
    const onComplete = (selectedType) => {
        handleResult(selectedType)

        closed.val = true
    }
    const onCancel = () => closed.val = true;
    const form = SelectBlockTypeForm({types, onComplete, onCancel})
    const modal = ModalBackground(ModalForeground({}, form))

    van.add(document.body, () => closed.val ? null : modal);
}

export {BlockRow, SelectBlockTypeModal};