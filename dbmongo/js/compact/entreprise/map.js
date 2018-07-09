function map() { 
    try{
        if (this.value != null) {
            emit(this.value.siren, this.value) 
        }           
    } catch(error) {
        print(this.value.siren)
    }
}