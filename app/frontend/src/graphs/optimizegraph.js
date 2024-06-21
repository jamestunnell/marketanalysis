import van from "vanjs-core"
import { nanoid } from 'nanoid'

import { Button, ButtonCancel, ButtonIcon, ButtonIconTooltip } from "../elements/buttons"
import { RangeIncl } from "../constraint"
import { IconCheck, IconClose, IconError } from "../elements/icons"
import { allMeasurements } from './measurement.js'
import { ModalBackground, ModalForeground } from "../modal"
import { IntRange, NumberRange } from "../elements/number"
import { allObjectives } from './objective.js'
import Textbox from "../elements/textbox"
import Select from "../elements/select"
import { TargetParamsTable } from './targetparam.js'

const {div, option, span} = van.tags

import { PostJSON } from '../backend.js'

function newRandomSeed() {
    return Math.floor(Math.random() * Number.MAX_SAFE_INTEGER) - Number.MAX_SAFE_INTEGER
}

const startOptimizeGraphJob = ({graph, jobID, symbol, days, sourceQuantity, targetParams, settings}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/optimize`
        const object = {graph, jobID, symbol, days, sourceQuantity, targetParams, settings}
        const options = {accept: 'application/json'}

        console.log("optimizing graph", object)

        PostJSON({route, object, options}).then(resp => {
            if (resp.status != 202) {
                resp.json().then(appErr => {
                    console.log("failed to run graph", appErr);
    
                    reject(appErr);    
                })
            }

            // Avoid Fetch failed loading
            resp.text().then(text => resolve())
        }).catch(err => {
            console.log("failed to make optimize graph request", err)
            
            reject({
                title: "Action Failed",
                message: "failed to make optimize graph request",
                details: [err.message],
            })
        });
    });
}

const OptimizeGraphModal = ({graph, symbolSetting, sourceAddresses, targetParams}) => {
    const symbol = van.state(symbolSetting.value.val)
    
    const days = van.state(30)
    const daysConstraint = new RangeIncl(10, 1000)
    const daysErr = van.derive(() => daysConstraint.validate(days.val))

    const maxIter = van.state(100)
    const maxIterConstraint = new RangeIncl(1, 10000)
    const maxIterErr = van.derive(() => maxIterConstraint.validate(maxIter.val))

    const seed = van.state(newRandomSeed())

    const source = van.state('')
    const sourceOpts = [ option({value:'', selected: true}, '') ]

    const measurement = van.state('')
    const measurementOpts = [ option({value: '', selected: true}, '') ]

    const objective = van.state('')
    const objectiveOpts = [ option({value: '', selected: true}, '') ]

    allObjectives().forEach(o => {
        objectiveOpts.push(option({value: o}, o))
    })

    allMeasurements().forEach(m => {
        measurementOpts.push(option({value: m}, m))
    })

    sourceAddresses.forEach(s => {
        sourceOpts.push(option({value: s}, s))
    })

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

        "Max Iterations",
        IntRange({
            constraint: maxIterConstraint,
            value: maxIter,
        }),
        ButtonIconTooltip({
            icon: () => maxIterErr.val ? IconError() : IconCheck(),
            text: () => maxIterErr.val ? `Value is invalid: ${maxIterErr.val.message}` : "Value is valid",
        }),

        "Random Seed",
        NumberRange({parse: parseInt, value: seed}),
        span(),

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
        ButtonIconTooltip({
            icon: () => measurement.val.length === 0 ? IconError() : IconCheck(),
            text: () => measurement.val.length === 0 ? "No measurement selected" : "Value is valid",
        }),

        "Objective",
        Select({
            onchange: (e) => objective.val = e.target.value,
            options: objectiveOpts,
        }),
        ButtonIconTooltip({
            icon: () => objective.val.length === 0 ? IconError() : IconCheck(),
            text: () => objective.val.length === 0 ? "No objective selected" : "Value is valid",
        }),
    )

    const closed = van.state(false)
    const closeBtn = ButtonIcon({
        icon: IconClose(),
        onclick: (e) => closed.val = true,
    })

    closeBtn.classList.add("self-end")

    const modal = ModalBackground(ModalForeground({}, div(
        { class: "flex flex-col"},
        div(
            {class: "flex flex-row-reverse p-2"},
            closeBtn,
        ),
        mainForm,
        TargetParamsTable(targetParams),
        div(
            {class: "flex flex-row-reverse"},
            Button({
                child: "Run",
                onclick: () => {
                    const opts = {
                        jobID: nanoid(),
                        graph,
                        symbol: symbol.val,
                        days: days.val,
                        sourceQuantity: {
                            address: source.val,
                            measurement: measurement.val,
                        },
                        targetParams: targetParams.filter(t => t.selected.val).map(t => {
                            return {address: t.address, min: t.min.val, max: t.max.val}
                        }),
                        settings: {
                            objective: objective.val,
                            algorithm: 'SimulatedAnnealing',
                            randomSeed: seed.val,
                            maxIterations: maxIter.val,
                            keepHistory: false,
                        },
                    }
                    
                    startOptimizeGraphJob(opts).then(() => {
                        console.log('graph optimization started')
                    }).catch(appErr => {
                        AppErrorModal(appErr)
                    })
                },
                disabled: van.derive(() => {
                    if (symbol.val.length === 0) {
                        return true
                    }

                    if (daysConstraint.validate(days.val)) {
                        return true
                    }

                    if (maxIterConstraint.validate(maxIter.val)) {
                        return true
                    }

                    if (source.val.length === 0) {
                        return true
                    }
                    
                    if (measurement.val.length === 0) {
                        return true
                    }

                    if (objective.val.length === 0) {
                        return true
                    }

                    if (targetParams.filter(t => t.err.val).length > 0) {
                        return true
                    }

                    return false
                })
            }),
            ButtonCancel({
                child: "Cancel",
                onclick: (e) => closed.val = true,
            })
        ),
    )))

    van.add(document.body, () => closed.val ? null : modal)
}

export { OptimizeGraphModal }