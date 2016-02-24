/// <reference path="../typings/tsd.d.ts" />

import * as React from 'react'
import {render} from 'react-dom'
import {Provider} from 'react-redux'
import {Store, createStore} from 'redux'
import reducer from './reducers/rootReducer'
import App from './containers/App'

const initialState = {}
const store: Store = createStore(reducer, initialState)

render(
  <Provider store={store}>
    <App/>
  </Provider>,
  document.getElementById('root')
)
