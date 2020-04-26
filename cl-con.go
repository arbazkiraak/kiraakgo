//  fetch response content length with concurrency
import (
  "sync"
)

var urls = []string{"https://google.com","https://facebook.com","https://airbnb.com"}

type Job struct {  
    id       int
    randomno string
}

type Result struct {  
    job         Job
    sumofdigits int64
}


var jobs = make(chan Job, 10)  
var results = make(chan Result, 10)

// remove slice element by index
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// int64 because `nbytes` is of int64 datatype
func url_process(url string) int64{
    start_time := time.Now()
    //fmt.Printf("[+] REQUEST TO %s\n",url)
    resp,err := http.Get(url)
    if err != nil{
        fmt.Print(err) 
    }
    nbytes,err := io.Copy(ioutil.Discard,resp.Body)
    resp.Body.Close()
    if err != nil {
        fmt.Printf("while reading %s: %v\n", url, err)
    }
    time_taken := time.Since(start_time).Seconds()
    
    return nbytes
}


func worker(wg *sync.WaitGroup) {  
    for job := range jobs {
        output := Result{job, url_process(job.randomno)}
        results <- output
    }
    wg.Done()
}

func createWorkerPool(noOfWorkers int) { 
    var wg sync.WaitGroup
    for i := 0; i < noOfWorkers; i++ {
        wg.Add(1)
        go worker(&wg)
    }
    wg.Wait()
    close(results)
}

func result(done chan bool) {  
    for result := range results {
        fmt.Printf("Job id: %d, url: %s , results: %d\n", result.job.id, result.job.randomno, result.sumofdigits)
    }
    done <- true
}

func allocate(noOfJobs int){
    for i:=0; i < noOfJobs; i++{
        randomurl := urls[0]
        urls := RemoveIndex(urls,0)
        job := Job{i,randomurl}
        jobs <- job
    }
    close(jobs)
}

func main(){
    start_time := time.Now()
    noOfJobs := len(urls)
    go allocate(noOfJobs)
    done := make(chan bool)
    go result(done)
    noOfWorkers := 100
    createWorkerPool(noOfWorkers)
    <- done
}
main()
