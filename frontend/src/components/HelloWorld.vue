<script setup>
defineProps({
  msg: {
    type: String,
    required: true,
  },
})

function btn() {
  console.log("123")
  window.go.gui.App.Print("123")
  alert("hhh")
}
</script>

<template>
  <div>
    <h1>{{ msg }}</h1>
    <h3>
      You’ve successfully created a project with
    </h3>
    <button @click="btn">wahhh</button>
  </div>
</template>

<style scoped></style>
