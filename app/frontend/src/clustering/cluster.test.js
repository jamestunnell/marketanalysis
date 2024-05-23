import { expect, test } from 'vitest'

import {kMeansAdaptive} from './cluster.js'

test('no datasets produces no clusters', () => {
    expect(kMeansAdaptive({})).toHaveLength(0)
})

test('one dataset produces one cluster', () => {
    const clusters = kMeansAdaptive({"a":[1,2,3]})
    
    expect(clusters).toHaveLength(1)

    expect(clusters[0]).toHaveLength(1)
    expect(clusters[0]).toContain("a")
})

test('three similar datasets produces one cluster', () => {
    const datasets = {
        "a": [1,3,4,3,2,5],
        "b": [2,4,4,5,1,3,3,4],
        "c": [1,3,2,2,4,4,5],
    }
    const clusters = kMeansAdaptive(datasets)
    
    expect(clusters).toHaveLength(1)

    expect(clusters[0]).toHaveLength(3)
    expect(clusters[0]).toContain("a")
    expect(clusters[0]).toContain("b")
    expect(clusters[0]).toContain("c")
})

test('two sets of similar data produce two clusters', () => {
    const datasets = {
        "a": [1,3,4,3,2,5],
        "b": [2,4,4,5,1,3,3,4],
        "c": [11,17,9,12,10,14,12],
        "d": [13,15,10,8,10,11,12],
    }
    const clusters = kMeansAdaptive(datasets)
    
    expect(clusters).toHaveLength(2)

    expect(clusters[0]).toHaveLength(2)
    expect(clusters[1]).toHaveLength(2)
})

test('three sets of similar data produce three clusters', () => {
    const datasets = {
        "a": [1,3,4,3,2,5],
        "b": [2,4,4,5,1,3,3,4],
        "c": [5,5,2,1,3,4,2],
        "d": [11,17,9,12,10,14,12],
        "e": [13,15,10,8,10,11,12],
        "f": [10,9,11,15,16,13,12,11],
        "g": [132, 98, 150, 140, 120, 110],
        "h": [99,84,130,110,123,131,108],
        "i": [111,98,120,140,155,130,125,119],
    }
    const clusters = kMeansAdaptive(datasets)
    
    expect(clusters).toHaveLength(3)

    expect(clusters[0]).toHaveLength(3)
    expect(clusters[1]).toHaveLength(3)
    expect(clusters[2]).toHaveLength(3)
})