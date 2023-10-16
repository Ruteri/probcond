var dag = {
  nodes: [],
  edges: [],
  questionnaire: {},
  given: []
};

window.dag = dag;

function addNode() {
  var nodeNameInput = document.getElementById('nodeName');
  var nodeNegationInput = document.getElementById('nodeNegation');
  if (nodeNameInput.value !== '') {
    var node = {
      value: nodeNameInput.value,
      negation: nodeNegationInput.value
    };
    dag.nodes.push(node);
    nodeNameInput.value = '';
    nodeNegationInput.value = '';
  }
  updateGraph();
}
function addEdge() {
  var srcNodeSelect = document.getElementById('srcNode');
  var dstNodeSelect = document.getElementById('dstNode');
  var srcNodeId = parseInt(srcNodeSelect.value);
  var dstNodeId = parseInt(dstNodeSelect.value);
  if (srcNodeId !== dstNodeId) {
    var srcNode = dag.nodes[srcNodeId].value;
    var dstNode = dag.nodes[dstNodeId].value;
    var edge = {
      src: srcNode,
      dst: dstNode
    };
    dag.edges.push(edge);
    srcNodeSelect.value = '';
    dstNodeSelect.value = '';
  }
  updateGraph();
}
function exportData() {
  var data = JSON.stringify(dag);
  var blob = new Blob([data], {
    type: 'application/json'
  });
  var fileName = prompt('Enter file name:', 'dag.json');
  var link = document.createElement('a');
  link.href = URL.createObjectURL(blob);
  link.download = fileName;
  link.click();
}
function importData(event) {
  var file = event.target.files[0];
  var reader = new FileReader();
  reader.onload = function (e) {
    var contents = e.target.result;
    var importedDAG = JSON.parse(contents);
    // Update the existing DAG variable
    dag = importedDAG;
    // Add any necessary code to update the UI or perform other operations with the imported data
    updateGraph();
    displayQuestionnaire(dag.questionnaire);
  };
  reader.readAsText(file);
}

