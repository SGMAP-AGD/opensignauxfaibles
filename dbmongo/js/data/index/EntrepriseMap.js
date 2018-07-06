function map() { 
        batches = Object.keys(this.value.batch).map(batch => batch)
        emit(this.value.siren, batches) 
}