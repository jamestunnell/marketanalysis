import van from "vanjs-core";

import {Button, ButtonCancel} from './elements/buttons.js';
import {ModalBackground, ModalForeground} from './modal.js';
import { AlertError } from "./alerts.js";

const {div, p} = van.tags

class WorkflowStep {
    constructor({title, makeElements, consumeInput, makePrev, makeNext}) {
        this.title = title;
        this.makeElements = makeElements;
        this.consumeInput = consumeInput;
        this.makePrev = makePrev
        this.makeNext = makeNext
    }
}

const StepForm = ({step, buttons}) => {
    return div(
        {class: "flex flex-col rounded-md"},
        p({class: "text-lg font-medium font-bold text-center"}, step.title),
        div({class: "flex flex-col space-y-4"}, step.makeElements()),
        div({class:"mt-4 flex justify-center"}, buttons),
    )
}

const StepModal = ({step, makePrev=null, makeNext=null, onCancel, onComplete}) => {
    const closed = van.state(false);
    let buttons = [];

    if (makePrev) {
        const onclick = () => {
            closed.val = true;

            van.add(document.body, makePrev());
        }

        buttons.push(Button({child: "Prev", onclick: onclick}))
    }

    const cancel = ButtonCancel({
        child: "Cancel",
        onclick: () => {
            closed.val = true;

            onCancel()
        },
    })
    const advance = Button({
        child: makeNext ? "Next" : "Complete",
        onclick: () => {
            const err = step.consumeInput()
            if (err) {
                alert = new AlertError({title: "Invalid Input", substance: err.message})

                van.add(document.body, alert.render());

                return
            }

            closed.val = true

            if (makeNext) {
                van.add(document.body, makeNext());
            } else {
                onComplete()
            }
        },
    })

    if (makeNext) {
        buttons.push(advance)
        buttons.push(cancel)
    } else {
        buttons.push(cancel)
        buttons.push(advance)
    }

    return () => closed.val ? null : ModalBackground(
        ModalForeground(
            {},
            StepForm({step, buttons})
        ),
    );
}

function makeStepModal({stepFns, index, onComplete, onCancel}) {
    return StepModal({
        step: stepFns[index](),
        makePrev: (index == 0) ? null : () => makeStepModal({stepFns, index: index-1, onComplete, onCancel}),
        makeNext: (index == (stepFns.length-1)) ? null : () => makeStepModal({stepFns, index: index+1, onComplete, onCancel}),
        onCancel: onCancel,
        onComplete: onComplete,
    });
}

const ModalWorkflow = ({stepFns, onComplete, onCancel}) => {
    if (stepFns.length === 0) {
        console.log("workflow is empty, returning")
        
        return
    }

    const modal = makeStepModal({stepFns: stepFns, index: 0, onComplete, onCancel});

    van.add(document.body, modal);
}

export {WorkflowStep, ModalWorkflow}