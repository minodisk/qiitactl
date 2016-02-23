/// <reference path="../typings/tsd.d.ts" />

import {combineReducers} from 'redux'
import {UPDATE_FILE_LIST} from './actions/files'

function fileList(state = [], action) {
  switch (action.type) {
    case UPDATE_FILE_LIST:
      return [
        ...state,
        {
          text: action.paths,
          complete: true,
        }
      ]
    default:
      return state
  }
}

const app = combineReducers({
  fileList: fileList,
})

export default app
