import van from "vanjs-core"
import { Tooltip } from 'vanjs-ui'

import { Get } from './backend.js'
import {Button, ButtonCancel, ButtonIcon, ButtonIconTooltip} from './buttons.js';
import capitalize from './capitalize.js';
import { IconCheck, IconCollapsed, IconDelete, IconError, IconExpanded } from "./icons.js";
import { ModalWorkflow, WorkflowStep } from "./workflow.js";
import { ParamValItem, validateParamVal } from './paramvals.js'
import { ModalBackground, ModalForeground } from "./modal.js";
import { TableRow } from './table.js';

const {div, input, label, li, ul, option, p, select, span} = van.tags

const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

function validateBlock({block, info, otherNames}) {
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
        this.deleted = van.state(false);

        this.type = block.type
        this.name = van.state(block.name)
        this.paramVals = Object.fromEntries(info.params.map(p => {
            return [p.name, van.state(block.paramVals[p.name] || p.default)]
        }))
        this.recordingFlags = info.outputs.map((o) => {
            return van.state(block.recording.indexOf(o.name) >= 0)
        });
    }

    getName() {
        return this.name.val
    }
    
    makeBlock() {
        return {
            name: this.name.val,
            type: this.type,
            paramVals: Object.fromEntries(Object.entries(this.paramVals).map(([name,value]) => [name, value.val])),
            recording: this.info.outputs.map((out,idx) => {
                return this.recordingFlags[idx].val ? out.name : ""
            }).filter(n => n !== "")
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
                this.parent.markChanged()
            },
        })
        const deleteBtn = ButtonIcon({
            icon: IconDelete(),
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
                block   : this.makeBlock(),
                info: this.info,
                otherNames: otherNames,
            })
        })
        const statusBtn = ButtonIconTooltip({
            icon: () => validateErr.val ? IconError() : IconCheck(),
            tooltipText: van.derive(() => validateErr.val ? `Block is invalid: ${validateErr.val.message}` : "Block is valid"),
        });

        const recordingListItems = this.info.outputs.map((out, i) => {
            const props = {
                id: out.name,
                type: "checkbox",
                onchange: e => {
                    this.recordingFlags[i].val = e.target.checked

                    this.parent.markChanged()
                },
            }
    
            if (this.recordingFlags[i].val) {
                props.checked = "checked"
            }
    
            return li(
                input(props, capitalize(out.name)),
                span(out.name),
            )
        });
        const paramListItems = this.info.params.map(p => {
            return ParamValItem({
                param: p,
                value: this.paramVals[p.name],
                onChange: (val) => this.parent.markChanged(),
            })
        });
        const paramList = ul({hidden: true}, paramListItems)
        const recordingList = ul({hidden: true}, recordingListItems)
        const expanded = van.state(false)
        const expandBtn = ButtonIcon({
            icon: () => expanded.val ? IconExpanded() : IconCollapsed(),
            onclick: () => {
                if (expanded.val) {
                    paramList.setAttribute("hidden", true)
                    recordingList.setAttribute("hidden", true)
                } else {
                    paramList.removeAttribute("hidden")
                    recordingList.removeAttribute("hidden")
                }

                expanded.val = !expanded.val
            },
        });
        const rowItems = [ expandBtn, nameInput, this.type,
            paramList, recordingList, deleteBtn, statusBtn]
    
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

const ConfigureBlockForm = ({info, block, otherNames, onComplete, onCancel}) => {
    console.log(`configuring block with type ${block.type}`)

    const name = van.state(block.name)
    const recordingFlags = info.outputs.map((o) => {
        return van.state(block.recording.indexOf(o.name) >= 0)
    });
    const paramVals = Object.fromEntries(info.params.map(p => {
        return [p.name, van.state(block.paramVals[p.name] || p.default)]
    }))
    
    const ok = Button({
        child: "OK",
        onclick: () => {
            if (name.val.length === 0) {
                console.log("Name is empty")
                return 
            } else if (otherNames.indexOf(name.val) >= 0) {
                console.log(`Name '${name.val}' is not unique`)
                return
            }

            const nonDefaultParamVals = {};

            info.params.forEach(p => {
                const value = paramVals[p.name].val

                if (value === p.default) {
                    return
                }

                const err = validateParamVal(p, value)
                if (err) {
                    console.log(`invalid value ${value} for param ${p.name}`, err)

                    return
                }

                nonDefaultParamVals[p.name] = value
            })

            const recordingNames = [];

            info.outputs.forEach((o, i) => {
                if (recordingFlags[i].val) {
                    recordingNames.push(o.name)
                }
            })

            const blk = {
                name: name.val,
                recording: recordingNames,
                type: block.type,
                paramVals: nonDefaultParamVals,
            }

            console.log(`completing block configuration with type ${block.type}`)

            onComplete(blk)
        }
    })
    const cancel = ButtonCancel({child: "Cancel", onclick: onCancel})
    const paramValItems = info.params.map(p => {
        return ParamValItem({param: p, value: paramVals[p.name], onChange: (val) => this.parent.markChanged()})
    })
    const recordingItems = info.outputs.map((out, i) => {
        const props = {
            id: out.name,
            type: "checkbox",
            onchange: e => {
                recordingFlags[i].val = e.target.checked

                this.parent.markChanged()
            },
        }

        if (recordingFlags[i].val) {
            props.checked = "checked"
        }

        return li(
            input(props, capitalize(out.name)),
            span(out.name),
        )
    });

    return div(
        {class: "flex flex-col"},
        div(
            label({for: "name"}, "Name"),
            input({
                id: "name",
                class: inputClass,
                type: "text",
                value: block.name,
                placeholder: "Non-empty, unique",
                oninput: e => name.val = e.target.value,
            }),
        ),
        div(
            p("Param Values"),
            ul(paramValItems)
        ),
        div(
            p("Output Recording"),
            ul(recordingItems)
        ),
        div({class:"mt-4 flex flew-row-reverse"}, ok, cancel),
    )
}

const ConfigureBlockModal = ({info, block, otherNames, handleResult}) => {
    console.log(`other names`, otherNames)

    const closed = van.state(false)
    const onComplete = (block) => {
        handleResult(block)

        closed.val = true
    }
    const onCancel = () => closed.val = true;
    const form = ConfigureBlockForm({info, block, otherNames, onComplete, onCancel})
    const modal = ModalBackground(ModalForeground({}, form))

    van.add(document.body, () => closed.val ? null : modal);
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

export {BlockRow, ConfigureBlockModal, SelectBlockTypeModal};