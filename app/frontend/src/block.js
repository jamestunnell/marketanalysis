import van from "vanjs-core"

import { Get } from './backend.js'
import capitalize from './capitalize.js';
import { ModalWorkflow, WorkflowStep } from "./workflow.js";
import { ParamValItem, validateParamVal } from './paramvals.js'

const {div, input, label, li, ul, option, select, span} = van.tags

const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

class BlockWorkflow {
    constructor({block, infoByType, existingNames, handleResult}) {
        console.log("making block workflow", infoByType)

        this.infoByType = infoByType
        this.otherNames = existingNames.filter(name => name !== block.name)
        this.handleResult = handleResult

        this.name = block.name;
        this.type = block.type;
        this.paramVals = block.paramVals;
        this.recording = block.recording;
    }

    start() {
        console.log("starting block workflow", this.infoByType)
        
        const onComplete = () => {
            this.handleResult({
                name: this.name,
                type: this.type,
                paramVals: this.paramVals,
                recording: this.recording,
            })

            console.log("completed block workflow")
        }
        const onCancel = () => {
            console.log("canceled block workflow")
        };
        const stepFns = [
            () => this.makeNameStep(),
            () => this.makeTypeStep(),
            () => this.makeParamValsStep(),
            () => this.makeRecordingStep(),
        ]
        
        ModalWorkflow({stepFns, onComplete, onCancel})
    }
    
    makeNameStep() {
        const name = van.state(this.name)

        return new WorkflowStep({
            title: "Name",
            makeElements: () => {
                return [
                    label({for: "name"}, "Name"),
                    input({
                        id: "name",
                        class: inputClass,
                        type: "text",
                        value: this.name,
                        placeholder: "Non-empty, unique",
                        oninput: e => name.val = e.target.value,
                    }),
                ]
            },
            consumeInput: () => {
                if (name.val.length === 0) {
                    return new Error("Name is empty")
                } else if (this.otherNames.indexOf(name.val) >= 0) {
                    return new Error(`Name '${name.val}' is not unique`)
                }

                this.name = name.val

                console.log("consumed name input: %s", name.val, this.infoByType)

                return null
            },
        });
    }

    makeTypeStep() {
        const type = van.state(this.type)

        return new WorkflowStep({
            title: "Type",
            makeElements: () => {
                console.log("making elements for type step", this)
    
                const options = Object.keys(this.infoByType).map((t,i) => {
                    const props = {value: t}
                    if (t === this.type) {
                        props.selected = "selected"
                    }

                    return option(props, t);
                })
    
                const selectType = select(
                    { id: "type", class: inputClass, onchange: (e) => type.val = e.target.value },
                    options,
                )
    
                return [selectType]
            },
            consumeInput: () => {
                if (type.val === "") {
                    new Error("not type selected")
                }
    
                this.type = type.val
    
                return null
            },
        });
    }

    makeParamValsStep() {
        const info = this.infoByType[this.type]

        console.log("making param vals step for type %s", this.type, info)

        const values = {}
        
        Object.entries(this.paramVals).forEach(([name,value]) => {
            if (info.params.find(p => p.name == name)) {
                console.log("param %s has existing value %s", name, value)
                
                values[name] = van.state(value)
            }
        })

        return new WorkflowStep({
            title: "Parameter Values",
            makeElements: () => {
                const items = info.params.map(p => ParamValItem(p, values))

                return [ ul(items) ]
            },
            consumeInput: () => {
                const paramVals = {}
                const errs = []

                console.log("consuming param vals input", values)

                info.params.forEach(p => {
                    const v = values[p.name]
                    if (!v) {
                        console.log("no value found for param %s", p.name)

                        return
                    }
                    
                    const err = validateParamVal(p, v.val)

                    if (err) {
                        errs.push(err)
                    } else {
                        paramVals[p.name] = v.val
                    }
                })

                if (errs.length > 0) {
                    return errs[0]
                }

                this.paramVals = paramVals;

                return null
            },
        })
    }

    makeRecordingStep() {
        const info = this.infoByType[this.type]
        const recordingFlags = info.outputs.map((o) => {
            return van.state(this.recording.indexOf(o.name) >= 0)
        });
        
        return new WorkflowStep({
            title: "Recorded Outputs",
            makeElements: () => {
                const items = info.outputs.map((out, i) => {
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

                return [ ul(items) ]
            },
            consumeInput: () => {
                const recording = [];

                info.outputs.forEach((o, i) => {
                    if (recordingFlags[i].val) {
                        recording.push(o.name)
                    }
                })

                this.recording = recording

                return null
            },
        })
    }
}

export {BlockWorkflow};