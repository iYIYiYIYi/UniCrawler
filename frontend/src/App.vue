<script setup>
import HelloWorld from './components/HelloWorld.vue'
</script>

<template>
  <HelloWorld/>
</template>

<style>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
html::-webkit-scrollbar{
  width:10px;
  height:10px;
  /**/
}
html::-webkit-scrollbar-track{
  background: rgba(43, 43, 43, 0.1);
  border-radius:1px;
}
html::-webkit-scrollbar-thumb{
  background: #231b9c;
}
html::-webkit-scrollbar-thumb:hover{
  background: rgb(18, 152, 214);
}
html::-webkit-scrollbar-corner{
  background: #179a16;
}

html {
  overflow-x: hidden;
}
</style>
