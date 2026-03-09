import { Controller, Get, Param } from '@nestjs/common';
import { ProductsService } from './products.service';
import type { ProductResponse } from './product.response';
import { toProductResponse } from './product.response';

@Controller('products')
export class ProductsController {
  constructor(private readonly productsService: ProductsService) {}

  @Get()
  findAll(): ProductResponse[] {
    return this.productsService.findAll().map(toProductResponse);
  }

  @Get(':id')
  findOne(@Param('id') id: string): ProductResponse {
    return toProductResponse(this.productsService.findById(id));
  }
}
