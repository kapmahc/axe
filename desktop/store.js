import {createStore, applyMiddleware} from 'redux'
import thunkMiddleware from 'redux-thunk'

const initialState = {
  currentUser: null,
  siteInfo: null
}

export const actionTypes = {
  REFRESH: 'ADD',
  TICK: 'TICK'
}

// REDUCERS
export const reducer = (state = exampleInitialState, action) => {
  switch (action.type) {
    case actionTypes.TICK:
      return Object.assign({}, state, {
        lastUpdate: action.ts,
        light: !!action.light
      })
    case actionTypes.ADD:
      return Object.assign({}, state, {
        count: state.count + 1
      })
    default:
      return state
  }
}

// ACTIONS
export const serverRenderClock = (isServer) => dispatch => {
  return dispatch({
    type: actionTypes.TICK,
    light: !isServer,
    ts: Date.now()
  })
}

export const startClock = () => dispatch => {
  return setInterval(() => dispatch({type: actionTypes.TICK, light: true, ts: Date.now()}), 800)
}

export const addCount = () => dispatch => {
  return dispatch({type: actionTypes.ADD})
}

export const initStore = (initialState = exampleInitialState) => {
  return createStore(reducer, initialState, composeWithDevTools(applyMiddleware(thunkMiddleware)))
}
