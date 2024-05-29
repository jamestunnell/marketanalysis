import skmeans from 'skmeans'

function kMeansAdaptive(namedDatasets) {
    const names = Object.keys(namedDatasets)
    const datasets = names.map(name => namedDatasets[name])

    console.log("clustering with datasets", names)

    if (names.length === 0) {
        return []
    }
    
    if (names.length === 1) {
        return [names[0]]
    }

    const allValues = []
    const dsIndices = []

    datasets.forEach((values, dsIdx) => {
        values.forEach(val => {
            if (val) {
                allValues.push(val)
                dsIndices.push(dsIdx)
            }
        })
    })

    let clusters = []
    let k = names.length

    for (; k > 1; k--) {
        const majority = Number(k) / Number(k+1)

        console.log(`trying clustering with k=${k}`)

        const results = skmeans(allValues, k)

        console.log(`clustering with k=${k} complete`, results)

        const votesByDS = Object.fromEntries(names.map(name => {
            const votes = {}
            for (let j = 0; j < k; j++) {
                votes[j] = 0
            }

            return [name, votes]
        }))
    
        results.idxs.forEach((clusterIdx, valueIdx) => {
            const name = names[dsIndices[valueIdx]]
            const votes = votesByDS[name]
    
            votes[clusterIdx]++
        })

        const clustersByDS = {}

        names.forEach((name,i) => {
            const votes = votesByDS[name]
            const n = datasets[i].length
            for (let j = 0; j < k; j++) {
                const count = votes[j]

                const perc = Number(count) / Number(n)

                // console.log(`datastore ${name} votes for ${k}: ${100.0*perc}%`)

                if (perc >= majority) {
                    console.log(`datastore ${name}: selected cluster ${j}`)

                    clustersByDS[name] = j
                    
                    break
                }
            }
        })

        if (Object.keys(clustersByDS).length === names.length) {
            console.log(`clustering done with k=${k}`)

            for (let j = 0; j < k; j++) {
                const members = []

                names.forEach(name => {
                    if (clustersByDS[name] === j) {
                        members.push(name)
                    }
                })

                clusters.push(members)
            }

            break
        }
    }

    if (k===1) {
        return [names]
    }

    return clusters.filter(cluster => cluster.length > 0)
}

export {kMeansAdaptive}
