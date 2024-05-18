import van from "vanjs-core"

import { Get } from './backend.js'
import { Button, ButtonCancel } from "./buttons.js";
import capitalize from './capitalize.js';
import { ModalWorkflow, WorkflowStep } from "./workflow.js";
import { ParamValItem, validateParamVal } from './paramvals.js'
import { ModalBackground, ModalForeground } from "./modal.js";

const {div, input, label, li, ul, option, p, select, span} = van.tags

const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

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
    const paramValItems = info.params.map(p => ParamValItem(p, paramVals[p.name]))
    const recordingItems = info.outputs.map((out, i) => {
        const props = {
            id: out.name,
            type: "checkbox",
            onchange: e => recordingFlags[i].val = e.target.checked,
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

export {ConfigureBlockModal, SelectBlockTypeModal};