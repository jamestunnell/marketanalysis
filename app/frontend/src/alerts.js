import van from "vanjs-core"

const {button, div, h3} = van.tags

import { ButtonIcon } from './buttons.js'
import { IconClose, IconError } from "./icons"

class AlertError {
    constructor({title="Error", substance, secondsLeft=30}) {
        this.title = title
        this.substance = substance
        this.secondsLeft = van.state(secondsLeft)
        this.alert = div(
            { class: "block min-w-[25%] max-w-[50%] flex flex-col space-y-3 p-6 text-red-800 border border-red-300 rounded-lg bg-red-50 dark:bg-gray-800 dark:text-red-400 dark:border-red-800" },
        ) 
    }

    onTick() {
        if (this.secondsLeft.val < 1) {
            this.alert.remove();

            return
        }

        this.secondsLeft.val--

        setTimeout(() => this.onTick(), 1000);
    }

    render() {
        const dismissBtn = button(
            {
                type:"button",
                class: "text-red-800 bg-transparent border border-red-800 hover:bg-red-900 hover:text-white focus:ring-4 focus:outline-none focus:ring-red-300 font-medium rounded-lg text-xs px-3 py-1.5 text-center dark:hover:bg-red-600 dark:border-red-600 dark:text-red-500 dark:hover:text-white dark:focus:ring-red-800",
                onclick: () => this.alert.remove(),
            },
            van.derive(() => {
                return `Dismiss (${this.secondsLeft.val} seconds left)`
            }),
        )
        const closeBtn = ButtonIcon({
            icon: IconClose(),
            // text: "Close",
            onclick: ()=> this.alert.remove()},
        )
    
        setTimeout(() => this.onTick(), 1000);
    
        closeBtn.classList.add("self-end")
    
        van.add(this.alert,
            div(
                {class: "grid grid-cols-2 items-center"},
                div(
                    {class: "flex flex-row space-x-2"},
                    IconError(),
                    h3({class:"text-lg font-medium"}, this.title),
                ),
                div(
                    {class: "flex flex-row-reverse p-2"},
                    closeBtn,
                )
            ),
            this.substance,
            dismissBtn,
        )
    
        return this.alert
    }
}

export {AlertError}