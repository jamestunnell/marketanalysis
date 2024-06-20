import van from "vanjs-core"

import { Button, ButtonCancel, ButtonIconTooltip } from "../elements/buttons"
import { RangeIncl } from "../constraint"
import { IconCheck, IconError } from "../elements/icons"
import { allMeasurements } from './measurement.js'
import { ModalBackground, ModalForeground } from "../modal"
import { IntRange } from "../elements/number"
import Textbox from "../elements/textbox"
import Select from "../elements/select"
import Checkbox from "../elements/checkbox.js"
import { Table, TableRow } from "../elements/table.js"

const {div, option, span} = van.tags

import { PostJSON } from '../backend.js'

const optimizeGraph = ({graph, symbol, days, sourceQuantity, targetParams, settings}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/optimize`
        const object = {graph, symbol, days, sourceQuantity, targetParams, settings}
        const options = {accept: 'application/json'}

        console.log("running graph", object)

        PostJSON({route, object, options}).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to run graph", appErr);
    
                    reject(appErr);    
                })
            }

            resp.json().then(obj => resolve(obj))
        }).catch(err => {
            console.log("failed to make run graph request", err)
            
            reject({
                title: "Action Failed",
                message: "failed to make run graph request",
                details: [err.message],
            })
        });
    });
}

const OptimizeGraphModal = ({graph, symbolSetting, sources, params}) => {
    const symbol = van.state(symbolSetting.value.val)
    
    const days = van.state(30)
    const daysConstraint = new RangeIncl(10, 1000)
    const daysErr = van.derive(() => daysConstraint.validate(days.val))
    
    const source = van.state('')
    const sourceOpts = [ option({value:'', selected: true}, '') ]

    const measurement = van.state('')
    const ms = allMeasurements()
    const measurementOpts = ms.map(m => {
        return option({value: m, selected: m === ms[0]}, m)
    })

    sources.forEach(s => sourceOpts.push(option({value: s}, s)))

    const mainForm = div(
        {class: "grid grid-cols-3"},

        "Symbol",
        Textbox({value: symbol}),
        ButtonIconTooltip({
            icon: () => symbol.val.length === 0 ? IconError() : IconCheck(),
            text: () => symbol.val.length === 0 ? 'Value is empty' : 'Value is valid',
        }),

        "Days",
        IntRange({
            constraint: daysConstraint,
            value: days,
        }),
        ButtonIconTooltip({
            icon: () => daysErr.val ? IconError() : IconCheck(),
            text: () => daysErr.val ? `Value is invalid: ${daysErr.val.message}` : "Value is valid",
        }),

        "Source",
        Select({
            onchange: (e) => source.val = e.target.value,
            options: sourceOpts,
        }),
        ButtonIconTooltip({
            icon: () => source.val.length === 0 ? IconError() : IconCheck(),
            text: () => source.val.length === 0 ? "No source selected" : "Value is valid",
        }),

        "Measurement",
        Select({
            onchange: (e) => measurement.val = e.target.value,
            options: measurementOpts,
        }),
        span(),
    )

    const paramFlags = params.map(param => van.state(false))
    const paramRows = params.map((param, idx) => {
        const checkbox = Checkbox({checked: paramFlags[idx]})
        
        return TableRow([param, checkbox])
    })
    const paramsTable = Table({
        columnNames: ["Parameter", "Optimize"],
        tableBody: tbody({class:"table-auto"}, paramRows),
    })
    
    div(
        params.map(p => Checkbox)
    )

    modal = ModalBackground(ModalForeground({}, div(
        { class: "flex flex-col"},
        mainForm,
        paramsTable,
        div(
            {class: "flex flex-row-reverse"},
            Button({
                child: "OK",
                onclick: () => {
                    const targetParams = []
                    
                    paramFlags.forEach((flag, idx) => {
                        if (flag.val) {
                            targetParams.push(params[idx])
                        }
                    })

                    const opts = {
                        graph,
                        symbol: symbol.val,
                        days: days.val,
                        sourceQuantity: {
                            address: source.val,
                            measurement: measurement.val,
                        },
                        targetParams,
                        settings: {
                            randomSeed: 1,
                            algorithm: 'SimulatedAnnealing',
                            maxIterations: 10000,
                            keepHistory: false,
                        },
                    }
                    
                    optimizeGraph(opts).then(result => {
                        console.log('optimize succeeded', {result})
                    }).catch(appErr => {
                        AppErrorModal(appErr)
                    })
                },
                disabled: van.derive(() => {
                    if (symbol.length === 0) {
                        return true
                    }

                    if (daysConstraint.validate(days.val)) {
                        return true
                    }

                    return false
                })
            }),
            ButtonCancel()
        ),
    )))
}

export { OptimizeGraphModal }