import { TestBed } from '@angular/core/testing';

import { RoomBackendService } from './roombackend.service';

describe('RoomBackendService', () => {
  let service: RoomBackendService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(RoomBackendService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
