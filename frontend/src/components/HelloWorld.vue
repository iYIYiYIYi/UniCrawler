<script setup>
import {onMounted, reactive, ref} from 'vue'
import {Greet, StartCrawler, GetAfterTime, GetAllSchools, GetWithKeyWords} from '../../wailsjs/go/main/App'

const data = reactive({
  name: "计算机,工科,电子,通信,408",
  resultText: "Please enter your name below 👇",
  schools: [],
  log_s: "无事发生",
  isCrawlerWorking: false,
  isSearching: false,
  counter: 0,
})

function greet() {
  if (data.name.trim() === "") {
    getAllSchools()
  } else {
    getWithKeyWords()
  }
}

function log(s) {
  data.log_s = s
}

function getAllSchools() {
  log("获取中...")
  data.isSearching = true
  GetAllSchools().then(result => {
    result.sort((a, b)=>{
      return b.PublishTime < a.PublishTime ? -1 : 1
    })
    data.schools = result
    console.log(result)
    log("获取完成")
    data.isSearching = false
  })
}

function getWithKeyWords() {
  let keywords = data.name.split(",")
  log(keywords + "获取中...")
  data.isSearching = true
  GetWithKeyWords(keywords).then(result=> {
    result.sort((a, b)=>{
      return b.PublishTime < a.PublishTime ? -1 : 1
    })
    data.schools = result
    log("获取完成")
    data.isSearching = false
  })
}

function startCrawler() {
  log("正在爬取...")
  data.isCrawlerWorking = true
  StartCrawler().then(result => {
    greet()
    data.isCrawlerWorking = false
    clearInterval(ti.value)
  })
}

const timer = ref(0)
onMounted(()=>{
  greet()

  timer.value = window.setInterval(()=>{
    if (data.isCrawlerWorking) {
      return
    }
    startCrawler()
  }, 1000 * 60 * 30) 
})

</script>

<template>
  <main>
    <div id="input" class="input-box">
      <input id="name" v-model="data.name" autocomplete="off" class="input" type="text"/>
      <button class="btn" @click="greet" :disabled="data.isSearching">搜索</button>
      <button class="btn" @click="startCrawler" :disabled="data.isCrawlerWorking">启动爬虫</button>
    </div>
    <div class="container" style="margin-bottom: 10px;">
      <div class="log" >
        <p>
          {{ data.log_s }}
        </p>
      </div>
      <div v-if="data.schools.length > 0" class="cards">
        <div v-for="school in data.schools" class="card">
          <div class="header">
            <a :href="school.Link">
              {{ school.Title }}
            </a>
          </div>
          <div class="content">
            <p style="font-size: 14px;white-space: pre-line;">
              {{ school.Detail }}
            </p>      
          </div>
          <div class="footer">
            <div>{{ school.Name }}</div>
            <div>{{ school.PublishTime }}</div>
          </div>
        </div>
      </div>
    </div>
    <div>
      再往下就没有了...
    </div>
  </main>
  <div class="loadingContainer" v-if="data.isCrawlerWorking">
    <div class="loadingdiv">
      <div class="loading"></div>
    </div>
    <div>加载中...</div>
  </div>
</template>

<style scoped>
div.loadingContainer {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  margin: auto;
  background-color: rgba(37, 37, 37, 0.5);
  width: 200px;
  height: 200px;
  vertical-align: middle;
  padding: 10px;
}

div.loadingContainer > img {
  width: 88%;
}

div.loadingdiv {
  width: 60%;
  height: 60%;
  margin: 20px auto;
}

.cards {
  width: 85vw;
  margin: 0 auto;
  text-align: left;
  margin-bottom: 10px;
}

div .card{
  width: 100%;
  background-color: #ffffff;
  box-shadow: rgba(37, 37, 37, 0.5) 2px 1px 5px;
  border-radius: 3px;
  margin-bottom: 10px;
}

div.header {
  padding: 5px 0 5px 10px;
  color: black;
  border-bottom: 1px rgb(228, 221, 221);
  border-bottom-style: solid;
}

div.header a{
  color: rgb(39, 39, 39);
  font-size: larger;
  font-weight: bolder;
  text-decoration: none;
}

div.content {
  color: rgb(49, 49, 49);
  padding: 0 2px 2px 10px;
}

div.footer {
  color: black;
  font-size: 12px;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  padding-left: 10px;
  padding-bottom: 5px;
  padding-right: 10px;
}

.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box .btn {
  width: 70px;
  height: 35px;
  line-height: 35px;
  border-radius: 0;
  border: none;
  padding: 0 8px;
  cursor: pointer;
}

.input-box .btn:hover {
  background: #7e27f0;
  color: #e7e2e2;
}

.input-box .input {
  border: none;
  border-radius: 0;
  outline: none;
  height: 35px;
  line-height: 35px;
  padding: 0 10px;
  background-color: rgba(240, 240, 240, 1);
  margin-top: 40px;
  width: 45vw;
  -webkit-font-smoothing: antialiased;
}

.input-box .input:hover {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.input-box .input:focus {
  border: none;
  background-color: rgba(255, 255, 255, 1);
}

.loading {
  width: 100%;
  height: 100%;
  border: 2px solid #000;
  border-top-color: transparent;
  border-radius: 100%;
  transform: translateZ(0);
  will-change: transform;
  animation: circle infinite 0.75s linear;
}

@keyframes circle {
  0% {
    transform: rotate(0);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
