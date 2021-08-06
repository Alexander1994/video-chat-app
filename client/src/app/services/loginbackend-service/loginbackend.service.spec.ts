import { TestBed } from '@angular/core/testing';

import { LoginBackendService } from './loginbackend.service';

describe('LoginBackendService', () => {
  let service: LoginBackendService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(LoginBackendService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
