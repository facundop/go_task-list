var tasksService = {
  fetch: function () {
    var todos = ""

    return todos
  },
  save: function (todos) {
    
  }
}

// app Vue instance
var app = new Vue({
  // app initial state
  data: {
    todos: tasksService.fetch()
  },

  methods: {
    addTodo: function () {
      
    },

    removeTodo: function (todo) {
     
    },

    editTodo: function (todo) {
     
    }
  },

})

app.$mount('.app')