import van from "vanjs-core"
import Datepicker from 'flowbite-datepicker/Datepicker'

import { INPUT_CLASS } from "./input"
import { Get, PostJSON, PutJSON } from "./backend"

const { div, input, label, nav, ul } = van.tags

const SETTINGS_CONTAINER_ID = "settings"

function loadSetting(name) {
    return new Promise((resolve, reject) => {
        console.log(`loading setting ${name}`)

        Get(`/settings/${name}`).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to load setting", appErr);

                    reject(appErr)
                })
        
                return;
            }

            resp.json().then(obj => {
                console.log("loaded setting", obj);

                resolve(obj)
            })
        }).catch(err => {
            reject({
                title: "Action Failed",
                message: "failed to make get request",
                details: [err.message],
            })
        })
    })
}


function storeSetting({name, value}) {
    return new Promise((resolve, reject) => {
        const object = {name, value}

        console.log("storing setting", object)

        PostJSON({route: "/settings", object}).then(postResp => {
            if (postResp.status === 204) {
                // Avoid Fetch failed loading
                postResp.text().then(text => {
                    resolve()
                })
            } else {
                PutJSON({route: `/settings/${name}`, object}).then(putResp => {
                    if (putResp.status != 204) {
                        putResp.json().then(appErr => {
                            console.log("failed to store setting", appErr);
    
                            reject(appErr)
                        })
                
                        return;
                    }

                    console.log("stored setting", object)
    
                    // Avoid Fetch failed loading
                    putResp.text().then(text => {
                        resolve()
                    })
                }).catch(err => {
                    reject({
                        title: "Action Failed",
                        message: `failed to make put request`,
                        details: [err.message],
                    })
                })
            }
        }).catch(err => {
            reject({
                title: "Action Failed",
                message: `failed to make post request`,
                details: [err.message],
            })
        })
    })
}

class GraphSettings {
    constructor() {
        this.date = van.state('')
        this.symbol = van.state('')

        this.symbolInput = input({
            id: "symbol",
            class: INPUT_CLASS,
            type: "text",
            placeholder: "SPY, QQQ, etc.",
            oninput: e => {
                this.symbol.val = e.target.value
            },
            onchange: e => storeSetting({name: "symbol", value: e.target.value}),
        })
        this.dateInput = input({
            id: "actionDate",
            class: INPUT_CLASS,
            type: "text",
            placeholder: 'Select date',
        });

        this.dateInput.addEventListener('changeDate', (e) => {
            console.log(`changed date to ${e.target.value}`, e)

            storeSetting({name: "date", value: e.target.value}).then(() => {
                this.date.val = e.target.value
            })
        })

        const datePickerOpts = {
            autohide: true,
            container: `#${SETTINGS_CONTAINER_ID}`,
            daysOfWeekDisabled: [0, 6], // disable saturday and sunday
            format: "yyyy-mm-dd",
            maxDate: new Date(), // today
            todayHighlight: true,
        }
        const datepicker = new Datepicker(this.dateInput, datePickerOpts)
    }

    load() {
        loadSetting("symbol").then(setting => {
            this.symbolInput.value = setting.value
            this.symbol.val = setting.value
        })

        loadSetting("date").then(setting => {
            this.dateInput.value = setting.value
            this.date.val = setting.value
        })
    }

    render() {
        return div(
            nav(
                {class: "nav font-semibold text-md bg-gray-100 shadow shadow-gray-200"},
                div(
                    {class: "container flex text-gray-300 w-full"},
                    ul(
                        {id: SETTINGS_CONTAINER_ID, class: "flex items-center font-semibold flex-wrap"},
                        label({for: "date"}, "Symbol"),
                        this.symbolInput,
                        label({for: "date"}, "Date"),
                        this.dateInput
    
                        // NavBarItem({text: 'Home', route: 'home', routeArgs: [], currentRoute}),
                        // NavBarItem({text: 'Securities', route: 'securities', routeArgs: [], currentRoute}),
                        // NavBarItem({text: 'Graphs', route: 'graphs', routeArgs: [], currentRoute}),
                    )
                )
            )
        ) 
    }
}

export default GraphSettings