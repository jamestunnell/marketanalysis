import van from "vanjs-core"

import { getSecurities } from "./securities"
import { INPUT_CLASS } from "./input"
import { Get, PostJSON, PutJSON } from "./backend"
import { AppErrorAlert } from './apperror'

const { div, nav, option, select, ul } = van.tags

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
        this.symbol = van.state('')
        this.selectSymbol = select({
            id: "symbol",
            class: INPUT_CLASS,
            oninput: e => {
                storeSetting({name: "symbol", value: e.target.value}).then(() => {
                    this.symbol.val = e.target.value
                })
            }
        }, option({value: ""}, ""))
    }

    load() {
        // clear all the existing options
        while (this.selectSymbol.firstChild) {
            this.selectSymbol.removeChild(this.selectSymbol.firstChild)
        }
        
        Promise.all([
            getSecurities(),
            loadSetting("symbol"),
        ]).then(values => {
            const symbols = values[0].map(obj => obj.symbol)
            const selectedSymbol = values[1].value
            const opts = symbols.map(sym => {
                return option({value: sym, selected: sym === selectedSymbol}, sym)
            });
    
            this.symbol.val = (symbols.indexOf(selectedSymbol) >= 0) ? selectedSymbol : ''
            
            van.add(this.selectSymbol, option({value: ""}, ""))
            van.add(this.selectSymbol, opts);
        }).catch(appErr => {
            console.log("failed to resolve all promises", appErr)
    
            AppErrorAlert(appErr)
        })
    }

    render() {
        return div(
            nav(
                {class: "nav font-semibold text-lg bg-gray-100 shadow shadow-gray-200"},
                div(
                    {class: "container flex text-gray-300 w-full"},
                    ul(
                        {class: "flex items-center font-semibold flex-wrap"},
                        "Settings",
                        this.selectSymbol,
    
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