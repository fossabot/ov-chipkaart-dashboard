## Plan

### V1

- [X] Create a way to authenticate using the API
- [X] Modify the API to get transactions within a data range.
- [X] Create the ability to download all transactions into mongodb
- [X] Create the ability to convert csv download using the same API
- [X] Create ability to store destination prices for NS in the cache (DB) (map with 2 way keys)\
- [X] Create stations code cache service
- [X] Create NS transactions filter
- [X] Create enrichment step for NS transactions
- [X] Create ability to store Enriched NS transactions
- [ ] Use dependency injection for components
- [ ] Create ability to process NS transactions without discount
- [ ] Create ability to process NS transactions (Dal Voordeel)
- [ ] Create ability to process NS transactions with 20% discount with weekends (Altijd voordeel)
- [ ] Create ability to filter all off peak NS transactions (Dal Vrij)
- [ ] Implement LFU cache with O(1) complexity


### V2
- [ ] Create ability to process RET transactions
- [ ] Create ability to process 20% 19 Euro transactions
- [ ] Create ability to process RET & non RET transactions
- [ ] Create ability to calculate fare when there's a 20% discount 

### Blogging

- [ ] LFU Cache implementation and use