import React from 'react'
import {Route} from 'react-router'

import IndexTags from './tags/Index'
import FormTag from './tags/Form'

const routes = [
  (< Route key = "forum.tags.edit" path = "/forum/tags/edit/:id" component = {
    FormTag
  } />),
  (< Route key = "forum.tags.new" path = "/forum/tags/new" component = {
    FormTag
  } />),
  (< Route key = "forum.tags.index" path = "/forum/tags" component = {
    IndexTags
  } />)
]

export default routes
