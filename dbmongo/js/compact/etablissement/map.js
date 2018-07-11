function map() { 
    try{
        if (this.value != null) {
            emit(this.value.siret, this.value) 
        }   
    } catch(error) {
        print(this.value.siret)
    }
}