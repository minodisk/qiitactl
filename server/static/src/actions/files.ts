import { createAction, Action } from 'redux-actions';
import { assign } from 'lodash';
import { Promise } from 'es6-promise';

import { File } from '../models/files';
import * as types from '../constants/ActionTypes';
import { Socket } from './socket'

const s = new Socket()

const requestFiles = createAction<void>(
  types.REQUEST_FILES
)

const recieveFiles = createAction<File>(
  types.RECIEVE_FILES,
  () => {
    console.log('recieveFiles')
    return null
  }
);

const fetchFiles = () => {
  console.log('fetchFiles')
  return (dispatch, getState) => {
    console.log(getState())
    dispatch(requestFiles())
    return Promise.resolve()
  }
}

const startFile = createAction<File>(
  types.START_FILE,
  (file: File) => file
);

export {
  fetchFiles,
  startFile
}
