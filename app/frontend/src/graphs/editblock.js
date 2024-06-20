import van from "vanjs-core"

import { Button, ButtonCancel, ButtonIcon, ButtonIconTooltip } from '../elements/buttons.js'
import { MakeConstraint } from '../constraint.js'
import { IconCheck, IconClose, IconDelete, IconError } from '../elements/icons.js'
import { InputsTable, MakeInputs } from './input.js'
import { ParamsTable, MakeParams } from './param.js'
import { OutputsTable, MakeOutputs } from './output.js'
import Textbox from '../elements/textbox.js'

const {div, span} = van.tags

function validateBlockConfig({config, info, otherNames}) {
    console.log("validating block config", config)

    if (config.name.length === 0) {
        return new Error("Name is empty")
    }

    if (otherNames.indexOf(config.name) >= 0) {
        return new Error(`Name '${config.name}' is not unique`)
    }

    for (let i = 0; i < config.parameters.length; i++) {
        const name = config.parameters[i].name
        const paramInfo = info.parameters.find(paramInfo => paramInfo.name === name)

        if (!paramInfo) {
            return new Error(`failed to find info for param ${config.parameters[i].name}`)
        }

        const value = config.parameters[i].value
        const constraint = MakeConstraint(paramInfo)
        const err = constraint.validate(value)

        if (err) {
            return new Error(`param ${param.name} has invalid value ${value}`)
        }        
    }

    for (let i = 0; i < config.inputs.length; i++) {
        if (!info.inputs.find(input => input.name === config.inputs[i].name)) {
            return new Error(`failed to find input ${config.inputs[i].name}`)
        }
    }
    
    for (let i = 0; i < config.outputs.length; i++) {
        if (!info.outputs.find(output => output.name === config.outputs[i].name)) {
            return new Error(`failed to find output ${config.outputs[i].name}`)
        }
    }

    return null
}

const TAB_NAME_PARAMS = "Params"
const TAB_NAME_INPUTS = "Inputs"
const TAB_NAME_OUTPUTS = "Outputs"

const EditBlockForm = ({config, info, otherNames, possibleSources, onComplete, onCancel, onDelete}) => {
    const type = config.type
    const nameWorking = van.state(config.name)
    const inputs = MakeInputs({infos: info.inputs, configs: config.inputs})
    const outputs = MakeOutputs({infos: info.outputs, configs: config.outputs})
    const params = MakeParams({infos: info.parameters, configs: config.parameters})
    const makeBlockConfig = () => {
        return {
            name: nameWorking.val,
            type: config.type,
            parameters: params.filter(p => !p.isValueDefault()).map(p => p.makeConfig()),
            inputs: inputs.filter(i => !i.isSourceEmpty()).map(i => i.makeConfig()),
            outputs: outputs.filter(o => !o.isMeasurementsEmpty()).map(o => o.makeConfig()),
        }
    }
    const validateErr = van.derive(() => {
        console.log("using other names", otherNames)
        
        return validateBlockConfig({
            config: makeBlockConfig(),
            info,
            otherNames,
        })
    })

    const currentTab = van.state(TAB_NAME_PARAMS)
    const inputsTable = InputsTable({
        inputs, possibleSources,
        hidden: van.derive(() => currentTab.val !== TAB_NAME_INPUTS),
    })
    const outputsTable = OutputsTable({outputs, hidden: van.derive(() => currentTab.val !== TAB_NAME_OUTPUTS)})
    const paramsTable = ParamsTable({params, hidden: van.derive(() => currentTab.val !== TAB_NAME_PARAMS)})
    
    const nameInput = Textbox({
        value: nameWorking,
        placeholder: "Non-empty, unique",
    })
    const closeBtn = ButtonIcon({icon: IconClose(), onclick: onCancel})
    const deleteBtn = ButtonIcon({icon: IconDelete(), onclick: onDelete});
    const statusBtn = ButtonIconTooltip({
        icon: () => validateErr.val ? IconError() : IconCheck(),
        text: van.derive(() => validateErr.val ? `Block is invalid: ${validateErr.val.message}` : "Block is valid"),
    });
    const ok = Button({
        child: "OK",
        disabled: validateErr,
        onclick: () => onComplete(makeBlockConfig()),
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
                {class: "flex flex-row"},
                Button({
                    child: TAB_NAME_PARAMS,
                    onclick: (e) => currentTab.val = TAB_NAME_PARAMS,
                    disabled: van.derive(() => currentTab.val === TAB_NAME_PARAMS),
                }),
                Button({
                    child: TAB_NAME_INPUTS,
                    onclick: (e) => currentTab.val = TAB_NAME_INPUTS,
                    disabled: van.derive(() => currentTab.val === TAB_NAME_INPUTS),
                }),
                Button({
                    child: TAB_NAME_OUTPUTS,
                    onclick: (e) => currentTab.val = TAB_NAME_OUTPUTS,
                    disabled: van.derive(() => currentTab.val === TAB_NAME_OUTPUTS),
                }),
            ),
            paramsTable,
            inputsTable,
            outputsTable,
        ),
        div({class: "flex flex-row-reverse"}, ok, cancel),
    )
}

export default EditBlockForm