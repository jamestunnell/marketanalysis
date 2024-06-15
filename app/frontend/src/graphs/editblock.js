import van from "vanjs-core"

import { Button, ButtonCancel, ButtonIcon, ButtonIconTooltip } from '../buttons.js';
import { IconCheck, IconClose, IconCollapsed, IconDelete, IconError, IconExpanded } from "../icons.js";
import { InputsTable } from './inputs.js'
import { ParamValsTable, validateParamVal } from './paramvals.js'
import { OutputsTable } from './outputs.js'

const {div, input, span} = van.tags

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

    for (let i = 0; i < block.inputs.length; i++) {
        if (!info.inputs.find(input => input.name === block.inputs[i].name)) {
            return new Error(`failed to find input ${block.inputs[i].name}`)
        }
    }
    
    for (let i = 0; i < block.outputs.length; i++) {
        if (!info.outputs.find(output => output.name === block.outputs[i].name)) {
            return new Error(`failed to find output ${block.outputs[i].name}`)
        }
    }

    return null
}

const EditBlockForm = ({block, info, otherNames, possibleSources, onComplete, onCancel, onDelete}) => {
    const type = block.type
    const nameWorking = van.state(block.name)
    const paramValsWorking = Object.fromEntries(info.params.map(p => {
        const nonDefaultVal = block.paramVals ? block.paramVals[p.name] : null

        return [p.name, van.state(nonDefaultVal ?? p.default)]
    }))
    const sourcesWorking = Object.fromEntries(info.inputs.map(input => {
        let source = ""
        
        if (block.inputs) {
            block.inputs.forEach(cfg => {
                if (cfg.name === input.name) {
                    source = cfg.source
                }
            })
        }

        return [input.name, van.state(source)]
    }))
    const measurementsWorking = Object.fromEntries(info.outputs.map(output => {
        let measurements = []
        
        if (block.outputs) {
            block.outputs.forEach(cfg => {
                if (cfg.name === output.name) {
                    measurements = cfg.measurements
                }
            })
        }

        return [output.name, van.state(measurements)]
    }))
    const paramsCollapsed = van.state(false)
    const inputsCollapsed = van.state(false)
    const outputsCollapsed = van.state(false)
    const modifiedBlock = van.derive(() => {
        const paramVals = {}
        const inputs = []
        const outputs = []

        Object.entries(paramValsWorking).forEach(([name, value]) => {
            if (value.val !== info.params.find(p => p.name === name).default) {
                paramVals[name] = value.val
            }
        })

        Object.entries(sourcesWorking).forEach(([name, source]) => {
            if (source.val.length > 0) {
                inputs.push({name, source: source.val})
            }
        })

        Object.entries(measurementsWorking).forEach(([name, measurements]) => {
            if (measurements.val.length > 0) {
                console.log(`pushing measurements for ${name}: ${measurements}`)
                
                outputs.push({name, measurements: measurements.val})
            }
        })

        return {
            name: nameWorking.val,
            type: block.type,
            paramVals, inputs, outputs,
        }
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
        onclick: () => onComplete(modifiedBlock.val),
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
                InputsTable({
                    inputs: info.inputs,
                    sources: sourcesWorking,
                    possibleSources,
                }),
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
                OutputsTable({
                    outputs: info.outputs,
                    measurements: measurementsWorking,
                })
            ),
            span({class: "mb-2"}),
        ),
        div({class: "flex flex-row-reverse"}, ok, cancel),
    )
}

export {EditBlockForm, validateBlock}