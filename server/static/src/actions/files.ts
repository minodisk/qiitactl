import { createAction, Action } from 'redux-actions';
import { assign } from 'lodash';
import { Promise } from 'es6-promise';

import { File } from '../models/files';
import * as types from '../constants/ActionTypes';
import { socket } from './socket'


const requestFile = createAction<void>(
  types.REQUEST_FILES
)

const recieveFile = createAction<File>(
  types.RECIEVE_FILES,
  (file: File) => file
);

export const fetchFiles = () => {
  return (dispatch, getState) => {
    dispatch(requestFile())
    return socket.call('GetAllFiles', null)
      .then(file => dispatch(recieveFile(file)))
  }
}

export const startFile = createAction<File>(
  types.START_FILE,
  (file: File) => file
);
