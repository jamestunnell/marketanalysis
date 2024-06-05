import van from "vanjs-core"

import { Route } from 'vanjs-router'
import GraphsPage from './src/graphs.js'
import GraphPage from './src/graph.js'
import NavBar from './src/navbar.js'
import { SecuritiesPage } from './src/securities.js'
import GraphSettings from "./src/graphsettings.js"

import './index.css';

const RouteHome = () => {
    return Route(
        { name: 'home' },
        NavBar({currentRoute: 'home'}),
    )
}

const RouteSecurities = () => {
    const symbol = van.state('');
    
    return Route(
        {
            name: 'securities',
            onFirst() {
            },
            onLoad(route) {
                if (route.args.length == 0) {
                    symbol.val = '';

                    return
                }

                symbol.val = route.args[0]
            }
        },
        () => {
            if (symbol.val === '') {
                return NavBar({currentRoute: 'sercurities'})
            }

            return NavBar({currentRoute: `sercurities/${symbol.val}`})
        },
        () => {
            if (symbol.val === '') {
                return SecuritiesPage()
            }

            // return SecurityPage(symbol.val)
        }
    )
}

const RouteGraphs = () => {
    const graphID = van.state('')
    const settings = new GraphSettings()
    
    return Route(
        {
            name: 'graphs',
            onFirst() {
            },
            onLoad(route) {
                settings.load()
                
                if (route.args.length == 0) {
                    graphID.val = '';

                    return
                }

                graphID.val = route.args[0]
            }
        },
        () => {
            if (graphID.val === '') {
                return NavBar({currentRoute: 'graphs'})
            }

            return NavBar({currentRoute: `graphs/${graphID.val}`})
        },
        settings.render(),
        () => {
            if (graphID.val === '') {
                return GraphsPage()
            }

            return GraphPage({id: graphID.val, settings})
        }
    )
}

van.add(document.body, RouteHome(), RouteSecurities(), RouteGraphs())
