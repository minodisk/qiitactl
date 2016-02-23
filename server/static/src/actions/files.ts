export const UPDATE_FILE_LIST = 'UPDATE_FILE_LIST'
export const SET_VISIBILITY_FILTER = 'SET_VISIBILITY_FILTER'

export function updateFileList(paths) {
  return {
    type: UPDATE_FILE_LIST,
    paths: paths,
  }
}

export function setVisibilityFilter(filter) {
  return { type: SET_VISIBILITY_FILTER, filter }
}
