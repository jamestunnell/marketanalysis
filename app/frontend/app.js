import van from "vanjs-core"
import { Route } from 'vanjs-router'
import Graphs from './src/graphs.js'
import Graph from './src/graph.js'
import NavBar from './src/navbar.js'

const { div } = van.tags

const App = () => {
    return div(
        Route({ name: 'home' },
            NavBar({currentRoute: 'home'}),
        ),
        () => {
            const graphID = van.state('');
            return Route(
                {
                    name: 'graphs',
                    onFirst() {

                    },
                    onLoad(route) {
                        if (route.args.length == 0) {
                            return
                        }

                        graphID.val = route.args[0]
                    }
                },
                NavBar({currentRoute: 'graphs'}),
                () => {
                    if (graphID.val === '') {
                        return Graphs()
                    }
                    
                    console.log(`graph ID set to ${graphID.val}`)

                    return Graph(graphID.val)
                }
            )
        },
    )
}

van.add(document.body, App())
