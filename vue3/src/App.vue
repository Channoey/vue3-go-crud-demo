<template>
  <div class="table-box">
    <div class="title">
    <!-- 标题 -->
    <p><h2>最简单的CRUD</h2></p>
  </div>
  <!-- 搜索 -->
  <div class="query-box">
    <el-input class="query-input" v-model="queryInput" placeholder="请输入名字搜索" @input="handleQueryName"/>
    <div class="btn-list">
      <el-button type="primary" @click="handleAdd">增加</el-button>
      <el-button type="danger" @click="handleDelList" v-if="multipleSelection.length > 0">多选删除</el-button>
    </div>


  </div>
  <el-table ref="multipleTableRef"
    :data="tableData"
    style="width: 100%"
    @selection-change="handleSelectionChange"
    border>
    <el-table-column type="selection" width="55" />
    <el-table-column prop="name" label="姓名" width="120" />
    <el-table-column prop="email" label="邮箱" width="120" />
    <el-table-column prop="phone" label="电话" width="120" />
    <el-table-column prop="address" label="地址" width="300" />
    <el-table-column fixed="right" label="Operations" width="120">
      <template #default="scope">
        <el-button link type="primary" size="small" @click="handleRowDel(scope.row)" style="color:#F56C6C">删除</el-button>
        <el-button link type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
      </template>
    </el-table-column>
  </el-table>
    <el-pagination
        background
        layout="prev, pager, next"
        style="display: flex;justify-content: center;margin-top: 10px"
        :page-count="total2"
        v-model:current-page="curPage"
        @current-change="handlePageChange"
    />
    <el-dialog v-model="dialogFormVisible" :title="dialogType === 'add' ? '新增' : '编辑'">
    <el-form :model="tableForm">
      <el-form-item label="姓名" :label-width="80">
        <el-input v-model="tableForm.name" autocomplete="off" />
      </el-form-item>
      <el-form-item label="邮箱" :label-width="80">
        <el-input v-model="tableForm.email" autocomplete="off" />
      </el-form-item>
      <el-form-item label="电话" :label-width="80">
        <el-input v-model="tableForm.phone" autocomplete="off" />
      </el-form-item>
      <el-form-item label="状态" :label-width="80">
        <el-input v-model="tableForm.state" autocomplete="off" />
      </el-form-item>
      <el-form-item label="地址" :label-width="80">
        <el-input v-model="tableForm.address" autocomplete="off" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button type="primary" @click="dialogConfirm">
          确认
        </el-button>
      </span>
    </template>
  </el-dialog>
  </div>
  
  
</template>

<script setup>
import {ref} from "vue";
import request from "./utils/request.js";
// 数据
let queryInput = $ref("")
let tableData = $ref([
  {
    id:"1",
    name:"hhh1",
    email:"123@qq.com",
    phone:"1212121221",
    state:"在职",
    address:"123"
  },
  {
    id:"2",
    name:"hhh2",
    email:"123@qq.com",
    phone:"1212121221",
    state:"在职",
    address:"123"
  },
  {
    id:"3",
    name:"hhh3",
    email:"123@qq.com",
    phone:"1212121221",
    state:"在职",
    address:"123"
  },
  {
    id:"4",
    name:"hhh4",
    email:"123@qq.com",
    phone:"1212121221",
    state:"在职",
    address:"123"
  },
])
let multipleSelection = $ref([])
let dialogFormVisible = $ref(false)
let tableForm = $ref({
  name:'张三',
  email:"123@qq.com",
  phone:"123",
  state:"在职",
  address:"广东省"
})
let dialogType = $ref('add')
let tableDataCopy = Object.assign(tableData)
let total2 = $ref(1000)
let curPage = $ref(1)

//删除行
const handleRowDel = async ({ID}) => {
  // let index = tableData.findIndex(item=>item.id === id)
  // tableData.splice(index,1)
  await request.delete('/delete/'+ID)
  await getTableData(curPage)
}

//获得列表数据
const getTableData = async (cur = 1) =>{
  let res = await request.get('/list',{
    pageSize: 1,
    pageNum: cur
  })
  tableData = res.list
  total2 = res.total
  curPage = res.pageNum
  return res
}
console.log()

//请求分页
const handlePageChange = (val) =>{
  getTableData(val)
}

//新增
const handleAdd = () =>{
  dialogType = "add"
  dialogFormVisible = true
  tableForm = {}
}

//点击确认
const dialogConfirm = async () =>{
  dialogFormVisible = false
  //1.判断新增还是编辑
  if(dialogType==="add"){
    //1.拿到数据
    //2.展示数据
    // tableData.push({
    //   id:(tableData.length + 1).toString(),
    //   ...tableForm
    // })

    await request.post("/add",{
      ...tableForm
    })
    await getTableData(curPage)


  }
  if(dialogType==="edit"){
    // //1.获取索引
    // let index = tableData.findIndex(item => item.id===tableForm.id)
    // //2.替换数据
    //     tableData[index] = tableForm

    await request.put(`/update/${tableForm.ID}`,{
      ...tableForm
    })
    await getTableData(curPage)
  }

}

//选中
const handleSelectionChange = (val) => {
  multipleSelection = []
  val.forEach(item=>{
    multipleSelection.push(item.id)
  })
}

//多选删除
const handleDelList = () =>{
  multipleSelection.forEach(id=>{
    handleRowDel({id})
  })
  multipleSelection = []
}

//编辑
const handleEdit = (row) => {
  dialogFormVisible = true
  dialogType = "edit"
  tableForm = {...row}
}

//搜索
const handleQueryName = async (val) => {
  // if(val.length>0){
  //   tableData = tableData.filter(item=>item.name.toLowerCase().match(val.toLowerCase()))
  // }else{
  //   tableData = tableDataCopy
  // }
  if(val.length>0){
    tableData = await request.get(`/list/${val}`)
  }else{
    getTableData(curPage)
  }

}
</script>

<style scoped>
.table-box{
  width: 800px;
  margin: 240px;
}

.title{
  text-align: center;
}

.query-box{
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.query-input{
  width: 200px;
}
</style>
