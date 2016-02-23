export const UPDATE_FILE_LIST = 'UPDATE_FILE_LIST'

export function updateFileList(paths) {
  return {
    type: UPDATE_FILE_LIST,
    paths: paths,
  }
}
